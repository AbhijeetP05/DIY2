package services

import (
	"DIY2/models"
	"DIY2/utils"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Orders struct {
	conn *gorm.DB
}

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
	orders, err := orderModel.GetTopOrders(o.conn)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, orders)
}

func (o *Orders) TopProductsForAllStores(w http.ResponseWriter, r *http.Request) {
	storeModel := models.StoreModel{}
	stores, err := storeModel.GetAllStores(o.conn)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	m := make(map[int64][]int64)
	c := make(chan map[int64][]int64, len(stores))
	for i := 0; i < len(stores); i++ {
		go func(c chan map[int64][]int64) {
			orderModel := models.OrderModel{StoreId: stores[i]}
			products, err := orderModel.GetTopOrders(o.conn)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}

			m[orderModel.StoreId] = products
			c <- m
		}(c)

		m = <-c
	}
	resp := make([]TopProductsResponse, len(stores))
	count := 0
	for key, value := range m {
		fmt.Println(key, value)
		if len(value) > 2 {
			value = value[0:2]
		}
		resp[count] = TopProductsResponse{StoreId: key, Products: value}
		count++
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func NewOrder(db *gorm.DB) *Orders {
	return &Orders{conn: db}
}
