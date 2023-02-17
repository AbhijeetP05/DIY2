package services

import (
	"DIY2/models"
	"errors"
	"fmt"
)

type Stores struct {
	storeRepo models.IStoreRepo
	orderRepo models.IOrderRepo
}

//go:generate mockgen -destination=../mocks/store_mock.go -package=mocks go-mux/services IStores
type IStores interface {
	GetProducts(id int64, limit, start int) ([]models.ProductModel, error)
	AddProducts(id int64, products []models.ProductModel) (map[string]string, error)
	BuyProduct(productId, storeId int64) (string, error)
}

func (s *Stores) GetProducts(id int64, limit, start int) ([]models.ProductModel, error) {

	store := models.StoreModel{StoreId: id}
	products := s.storeRepo.GetProductsInStore(&store, limit, start)

	if products == nil {
		return nil, errors.New("some error occurred")
	}
	return products, nil
}

func (s *Stores) AddProducts(id int64, products []models.ProductModel) (map[string]string, error) {

	store := models.StoreModel{StoreId: int64(id)}

	fmt.Printf("%v", products)
	result := s.storeRepo.AddProducts(&store, products)
	if !result {
		return nil, errors.New("some error occurred")
	}
	fmt.Println("Success")
	payload := map[string]string{"result": "success"}
	return payload, nil
}

func (s *Stores) BuyProduct(productId, storeId int64) (string, error) {

	storeModel := models.StoreModel{StoreId: storeId, ProductId: productId, IsAvailable: true}
	err := s.storeRepo.ProductExists(&storeModel)
	if err != nil {
		return "", err
	}
	orderModel := models.OrderModel{ProductId: storeModel.ProductId, StoreId: storeModel.StoreId}
	err = s.orderRepo.BuyProduct(&orderModel)
	if err != nil {
		return "", err
	}

	payload := fmt.Sprintf("{orderId: %v}", orderModel.Id)
	return payload, nil
}

func NewStore(orderRepo models.IOrderRepo, storeRepo models.IStoreRepo) *Stores {
	return &Stores{orderRepo: orderRepo, storeRepo: storeRepo}
}
