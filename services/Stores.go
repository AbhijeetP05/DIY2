package services

import (
	"DIY2/models"
	"DIY2/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Stores struct {
	conn *gorm.DB
}

//go:generate mockgen -destination=../mocks/store_mock.go -package=mocks go-mux/services IStores
type IStores interface {
	GetProducts(w http.ResponseWriter, r *http.Request)
	AddProducts(w http.ResponseWriter, r *http.Request)
	BuyProduct(w http.ResponseWriter, r *http.Request)
}

func (s *Stores) GetProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}
	store := models.StoreModel{StoreId: int64(id)}
	products := store.GetProductsInStore(s.conn, limit, start)

	if products == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Some Error Occurred")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, products)
}

func (s *Stores) AddProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	store := models.StoreModel{StoreId: int64(id)}
	var products []models.ProductModel
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&products); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Printf("%v", products)
	result := store.AddProducts(s.conn, products)
	if !result {
		utils.RespondWithError(w, http.StatusInternalServerError, "Some Error Occurred")
		return
	}
	fmt.Println("Success")
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (s *Stores) BuyProduct(w http.ResponseWriter, r *http.Request) {
	productId, _ := strconv.ParseInt(r.FormValue("productId"), 10, 64)
	storeId, _ := strconv.ParseInt(r.FormValue("storeId"), 10, 64)
	if productId == 0 || storeId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store or Product ID")
		return
	}
	storeModel := models.StoreModel{StoreId: storeId, ProductId: productId, IsAvailable: true}
	orderId, err := storeModel.BuyProduct(s.conn)
	if err != nil {
		if err.Error() == "record not found" {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	payload := fmt.Sprintf("{orderId: %v}", orderId)
	utils.RespondWithJSON(w, http.StatusOK, payload)
}

func NewStore(conn *gorm.DB) *Stores {
	return &Stores{conn: conn}
}
