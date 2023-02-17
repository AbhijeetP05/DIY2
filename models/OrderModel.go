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

type OrderRepo struct {
	conn *gorm.DB
}

//go:generate mockgen -destination=../mocks/orderModel_mock.go -package=mocks DIY2/models IOrderRepo
type IOrderRepo interface {
	BuyProduct(orderModel *OrderModel) error
	GetTopOrders(orderModel *OrderModel, interval int, limit int) ([]int64, error)
}

func NewOrderRepo(conn *gorm.DB) *OrderRepo {
	return &OrderRepo{conn: conn}
}

func (m *OrderRepo) BuyProduct(orderModel *OrderModel) error {
	result := m.conn.Create(&orderModel)
	return result.Error
}

func (m *OrderRepo) GetTopOrders(orderModel *OrderModel, interval int, limit int) ([]int64, error) {
	var products []int64
	whereQuery := fmt.Sprintf("sold_at>=now()-interval '%v hour' and STORE_ID = %v", interval, orderModel.StoreId)
	err := m.conn.Select("product_id").Group("product_id").Where(whereQuery).Order("count(product_id) desc").Limit(limit).Table("orders").Find(&products).Error
	return products, err
}

func (*OrderModel) TableName() string {
	return "orders"
}
