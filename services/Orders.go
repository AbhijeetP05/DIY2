package services

import (
	"DIY2/models"
)

type Orders struct {
	orderRepo models.IOrderRepo
	storeRepo models.IStoreRepo
}

//go:generate mockgen -destination=../mocks/order_mock.go -package=mocks DIY2/services IOrders
type IOrders interface {
	TopProductsInStore(storeId int64) ([]int64, error)
	TopProductsForAllStores() ([]TopProductsResponse, error)
}

type TopProductsResponse struct {
	StoreId  int64   `json:"store_id"`
	Products []int64 `json:"products"`
}

func (o *Orders) TopProductsInStore(storeId int64) ([]int64, error) {

	orderModel := models.OrderModel{StoreId: storeId}
	orders, err := o.orderRepo.GetTopOrders(&orderModel, 1, 5)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Orders) TopProductsForAllStores() ([]TopProductsResponse, error) {
	storeModel := models.StoreModel{}
	stores, err := o.storeRepo.GetAllStores(&storeModel)
	if err != nil {
		return nil, err
	}
	responseMap := make(map[int64][]int64)
	channel := make(chan map[int64][]int64, len(stores))
	var goErr error
	for i := 0; i < len(stores); i++ {
		go func(c chan map[int64][]int64) {
			err := o.getTopOrdersAllStores(c, stores[i], responseMap)
			if err != nil {
				goErr = err
				return
			}
		}(channel)
		if goErr != nil {
			return nil, goErr
		}
		responseMap = <-channel
	}
	payload := make([]TopProductsResponse, len(stores))
	count := 0
	for key, value := range responseMap {
		payload[count] = TopProductsResponse{StoreId: key, Products: value}
		count++
	}
	return payload, nil
}

func (o *Orders) getTopOrdersAllStores(c chan map[int64][]int64, storeId int64, responseMap map[int64][]int64) error {
	orderModel := models.OrderModel{StoreId: storeId}
	products, err := o.orderRepo.GetTopOrders(&orderModel, 1, 2)
	if err != nil {
		return err
	}

	responseMap[orderModel.StoreId] = products
	c <- responseMap
	return nil
}

func NewOrder(storeRepo models.IStoreRepo, orderRepo models.IOrderRepo) *Orders {
	return &Orders{storeRepo: storeRepo, orderRepo: orderRepo}
}
