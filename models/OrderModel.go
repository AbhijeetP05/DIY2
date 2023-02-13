package models

import (
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

func (m *OrderModel) CreateOrder(db *gorm.DB) error {
	result := db.Create(&m)
	return result.Error
}

func (m *OrderModel) GetTopOrders(db *gorm.DB) ([]int64, error) {
	var orders []OrderModel

	//result := db.Raw("select product_id from (select * from orders where SOLD_AT >= now() - interval '1 hour' and STORE_ID = ?) as \"o\" group by PRODUCT_ID order by count(PRODUCT_ID) desc limit 5", m.StoreId).Scan(&orders)
	result := db.Select("product_id").Group("product_id").Where("sold_at>=now()-interval '1 hour' and STORE_ID = ?", m.StoreId).Order("count(product_id) desc").Limit(5).Table("orders").Find(&orders)

	var products []int64
	for i := 0; i < len(orders); i++ {
		products = append(products, orders[i].ProductId)
	}
	return products, result.Error
}

func (*OrderModel) TableName() string {
	return "orders"
}
