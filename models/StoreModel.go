package models

import (
	"fmt"
	"gorm.io/gorm"
)

type StoreModel struct {
	StoreId     int64        `json:"storeId" gorm:"primaryKey;autoincrement:false"`
	ProductId   int64        `json:"productId" gorm:"primaryKey;autoincrement:false"`
	Product     ProductModel `gorm:"foreignKey:ProductId;references:ID"`
	IsAvailable bool         `json:"isAvailable"`
}

type StoreRepo struct {
	conn *gorm.DB
}

//go:generate mockgen -destination=../mocks/storeModel_mock.go -package=mocks DIY2/models IStoreRepo
type IStoreRepo interface {
	GetProductsInStore(storeModel *StoreModel, limit, start int) []ProductModel
	AddProducts(storeModel *StoreModel, products []ProductModel) bool
	ProductExists(storeModel *StoreModel) error
	GetAllStores(storeModel *StoreModel) ([]int64, error)
}

func NewStoreRepo(conn *gorm.DB) *StoreRepo {
	return &StoreRepo{conn: conn}
}

func (sr *StoreRepo) GetProductsInStore(s *StoreModel, limit, start int) []ProductModel {
	var productsInStore []StoreModel
	var products []ProductModel
	result := sr.conn.Model(&StoreModel{}).Where("store_id = ?", s.StoreId).Limit(limit).Offset(start).Find(&productsInStore)

	if result.Error != nil {
		fmt.Println("Some error occurred")
		return nil
	}
	tx := sr.conn.Begin()
	for i := 0; i < len(productsInStore); i++ {
		p := ProductModel{ID: &productsInStore[i].ProductId}
		result := sr.conn.First(&p)
		if result.Error != nil {
			tx.Rollback()
			break
		}
		products = append(products, p)

	}
	if tx.Error != nil {
		return nil
	}
	tx.Commit()
	return products
}

func (sr *StoreRepo) AddProducts(s *StoreModel, products []ProductModel) bool {

	tx := sr.conn.Begin()
	for i := 0; i < len(products); i++ {
		res := sr.conn.Create(&products[i])
		fmt.Println(products[i].ID, products[i].Name)
		if res.Error != nil {
			tx.Rollback()
			break
		}

		s.ProductId = *products[i].ID
		s.IsAvailable = true
		result := sr.conn.Model(&s).Where("store_id = ? and product_id = ?", s.StoreId, s.ProductId).Updates(&s)
		if result.RowsAffected == 0 {
			result = sr.conn.Create(&s)
		}
		if result.Error != nil {
			tx.Rollback()
			break
		}
	}
	if tx.Error != nil {
		return false
	}
	err := tx.Commit().Error
	return err == nil
}

func (sr *StoreRepo) ProductExists(s *StoreModel) error {
	err := sr.conn.First(&s).Error
	return err
}

func (sr *StoreRepo) GetAllStores(s *StoreModel) ([]int64, error) {
	var stores []int64
	err := sr.conn.Model(&s).Select("distinct store_id").Order("store_id").Scan(&stores).Error

	return stores, err
}

func (s *StoreModel) TableName() string {
	return "stores"
}
