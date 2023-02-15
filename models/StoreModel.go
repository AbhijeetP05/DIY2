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

//go:generate mockgen -destination=../mocks/storeModel_mock.go -package=mocks go-mux/models IStoreModel
type IStoreModel interface {
	GetProductsInStore(db *gorm.DB, limit, start int) []ProductModel
	AddProducts(db *gorm.DB, products []ProductModel) bool
	BuyProduct(db *gorm.DB, productId int64) (int64, error)
}

func NewStore(storeId int64, productId int64, IsAvailable bool) *StoreModel {
	return &StoreModel{StoreId: storeId, ProductId: productId, IsAvailable: IsAvailable}
}

func (s *StoreModel) GetProductsInStore(db *gorm.DB, limit, start int) []ProductModel {
	var productsInStore []StoreModel
	var products []ProductModel
	result := db.Model(&StoreModel{}).Where("store_id = ?", s.StoreId).Limit(limit).Offset(start).Find(&productsInStore)

	if result.Error != nil {
		fmt.Println("Some error occurred")
		return nil
	}
	tx := db.Begin()
	for i := 0; i < len(productsInStore); i++ {
		p := ProductModel{ID: &productsInStore[i].ProductId}
		result := db.First(&p)
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

func (s *StoreModel) AddProducts(db *gorm.DB, products []ProductModel) bool {

	tx := db.Begin()
	for i := 0; i < len(products); i++ {
		res := db.Create(&products[i])
		fmt.Println(products[i].ID, products[i].Name)
		if res.Error != nil {
			tx.Rollback()
			break
		}

		s.ProductId = *products[i].ID
		s.IsAvailable = true
		result := db.Model(&s).Where("store_id = ? and product_id = ?", s.StoreId, s.ProductId).Updates(&s)
		if result.RowsAffected == 0 {
			result = db.Create(&s)
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

func (s *StoreModel) ProductExists(db *gorm.DB) error {
	err := db.First(&s).Error
	return err
}

func (s *StoreModel) GetAllStores(db *gorm.DB) ([]int64, error) {
	var stores []int64
	err := db.Model(&s).Select("distinct store_id").Order("store_id").Scan(&stores).Error

	return stores, err
}

func (s *StoreModel) TableName() string {
	return "stores"
}
