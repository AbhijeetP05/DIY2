package services

import (
	"DIY2/models"
	"errors"
	"fmt"
)

type Stores struct {
	storeRepo   models.IStoreRepo
	productRepo models.IProductRepo
}

//go:generate mockgen -destination=../mocks/store_mock.go -package=mocks DIY2/services IStores
type IStores interface {
	GetProducts(id int64, limit, start int) ([]models.ProductModel, error)
	AddProducts(id int64, products []models.ProductModel) (map[string]string, error)
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

func NewStore(productRepo models.IProductRepo, storeRepo models.IStoreRepo) *Stores {
	return &Stores{productRepo: productRepo, storeRepo: storeRepo}
}
