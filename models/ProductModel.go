package models

import (
	"fmt"
	"gorm.io/gorm"
)

type ProductModel struct {
	ID    *int64  `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

//go:generate mockgen -destination=../mocks/store_mock.go -package=mocks go-mux/services IStores
type IProductModel interface {
	GetProduct(db *gorm.DB) *gorm.DB
	GetProducts(db *gorm.DB, limit, start int) ([]ProductModel, *gorm.DB)
	CreateProduct(db *gorm.DB) *gorm.DB
	UpdateProduct(db *gorm.DB, newProduct *ProductModel) *gorm.DB
	DeleteProduct(db *gorm.DB) *gorm.DB
}

func (*ProductModel) TableName() string {
	return "products"
}

func (p *ProductModel) GetProduct(db *gorm.DB) *gorm.DB {
	result := db.First(&p)
	return result
}

func (p *ProductModel) GetProducts(db *gorm.DB, limit, start int) ([]ProductModel, *gorm.DB) {
	var products []ProductModel
	result := db.Model(ProductModel{}).Offset(start).Limit(limit).Find(&products)

	return products, result
}

func (p *ProductModel) CreateProduct(db *gorm.DB) *gorm.DB {
	fmt.Println(p.ID)
	result := db.Create(&p)

	return result
}

func (p *ProductModel) UpdateProduct(db *gorm.DB, newProduct *ProductModel) *gorm.DB {
	result := db.Model(&p).Updates(newProduct)

	return result
}

func (p *ProductModel) DeleteProduct(db *gorm.DB) *gorm.DB {
	result := db.Delete(&p)

	return result
}
