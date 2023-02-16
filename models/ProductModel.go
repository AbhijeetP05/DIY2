package models

import (
	"gorm.io/gorm"
)

type ProductModel struct {
	ID    *int64  `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductRepo struct {
	conn *gorm.DB
}

//go:generate mockgen -destination=../mocks/store_mock.go -package=mocks go-mux/services IStores
type IProductRepo interface {
	GetProduct(p *ProductModel) *gorm.DB
	GetProducts(p *ProductModel, limit, start int) ([]ProductModel, *gorm.DB)
	CreateProduct(p *ProductModel) *gorm.DB
	UpdateProduct(p *ProductModel, newProduct *ProductModel) *gorm.DB
	DeleteProduct(p *ProductModel) *gorm.DB
}

func NewProductRepo(conn *gorm.DB) *ProductRepo {
	return &ProductRepo{conn: conn}
}

func (*ProductModel) TableName() string {
	return "products"
}

func (pr *ProductRepo) GetProduct(p *ProductModel) *gorm.DB {
	result := pr.conn.First(&p)
	return result
}

func (pr *ProductRepo) GetProducts(p *ProductModel, limit, start int) ([]ProductModel, *gorm.DB) {
	var products []ProductModel
	result := pr.conn.Model(ProductModel{}).Offset(start).Limit(limit).Find(&products)

	return products, result
}

func (pr *ProductRepo) CreateProduct(p *ProductModel) *gorm.DB {
	//fmt.Println(p.ID)
	result := pr.conn.Create(&p)

	return result
}

func (pr *ProductRepo) UpdateProduct(p *ProductModel, newProduct *ProductModel) *gorm.DB {
	result := pr.conn.Model(&p).Updates(newProduct)

	return result
}

func (pr *ProductRepo) DeleteProduct(p *ProductModel) *gorm.DB {
	result := pr.conn.Delete(&p)

	return result
}
