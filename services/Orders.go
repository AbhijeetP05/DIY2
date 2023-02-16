package services

import (
	"DIY2/models"
	"DIY2/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Orders struct {
	orderRepo models.IOrderRepo
	storeRepo models.IStoreRepo
}

//go:generate mockgen -destination=../mocks/order_mock.go -package=mocks DIY2/services IOrders
type IOrders interface {
	TopProductsInStore(w http.ResponseWriter, r *http.Request)
	TopProductsForAllStores(w http.ResponseWriter, r *http.Request)
}

type TopProductsResponse struct {
	StoreId  int64   `json:"store_id"`
	Products []int64 `json:"products"`
}

func (o *Orders) TopProductsInStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeId, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}
	orderModel := models.OrderModel{StoreId: storeId}
	orders, err := o.orderRepo.GetTopOrders(&orderModel, 1, 5)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, orders)
}

func (o *Orders) TopProductsForAllStores(w http.ResponseWriter, r *http.Request) {
	storeModel := models.StoreModel{}
	stores, err := o.storeRepo.GetAllStores(&storeModel)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseMap := make(map[int64][]int64)
	channel := make(chan map[int64][]int64, len(stores))
	for i := 0; i < len(stores); i++ {
		go func(c chan map[int64][]int64) {
			err := o.getTopOrdersAllStores(c, stores[i], responseMap)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		}(channel)

		responseMap = <-channel
	}
	payload := make([]TopProductsResponse, len(stores))
	count := 0
	for key, value := range responseMap {
		payload[count] = TopProductsResponse{StoreId: key, Products: value}
		count++
	}
	utils.RespondWithJSON(w, http.StatusOK, payload)
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
