package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type OrderModel struct {
	Id        int64     `json:"id" gorm:"primary key"`
	StoreId   int64     `json:"storeId"`
	ProductId int64     `json:"productId"`
	SoldAt    time.Time `json:"soldAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type IOrderModel interface {
	CreateOrder(db *gorm.DB) error
	GetTopOrders(db *gorm.DB) ([]int64, error)
}

func (m *OrderModel) BuyProduct(db *gorm.DB) error {
	result := db.Create(&m)
	return result.Error
}

func (m *OrderModel) GetTopOrders(db *gorm.DB, interval int64, limit int) ([]int64, error) {
	var products []int64
	fmt.Println(interval)
	whereQuery := fmt.Sprintf("sold_at>=now()-interval '%v hour' and STORE_ID = %v", interval, m.StoreId)
	err := db.Select("product_id").Group("product_id").Where(whereQuery).Order("count(product_id) desc").Limit(limit).Table("orders").Find(&products).Error
	return products, err
}

func (*OrderModel) TableName() string {
	return "orders"
}
