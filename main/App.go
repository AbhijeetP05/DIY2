package main

import (
	"DIY2/models"
	"DIY2/services"
	"DIY2/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// App main app for the program to run, contains a gorm database instance and a router instance for routing
type App struct {
	DB          *gorm.DB
	Router      *mux.Router
	storeRepo   models.IStoreRepo
	orderRepo   models.IOrderRepo
	productRepo models.IProductRepo
	products    services.IProducts
	stores      services.IStores
	orders      services.IOrders
}

// Initialize This function initializes the given application which will initialize the database and routes
func (a *App) Initialize(host, port, username, password, dbname string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", host, username, password, dbname, port)
	var err error
	a.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("UNABLE TO CONNECT TO DATABASE.")
	}
	log.Println("Database Initialized")
	a.Router = mux.NewRouter()

	a.productRepo = models.NewProductRepo(a.DB)
	a.orderRepo = models.NewOrderRepo(a.DB)
	a.storeRepo = models.NewStoreRepo(a.DB)

	a.products = services.NewProduct(a.productRepo)
	a.orders = services.NewOrder(a.storeRepo, a.orderRepo)
	a.stores = services.NewStore(a.orderRepo, a.storeRepo)

	a.InitializeRoutes()
	log.Println("Routes Initialized")
}

func (a *App) Run(host, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
func (a *App) home(w http.ResponseWriter, r *http.Request) {
	j := "{services: not available}"
	res, err := json.Marshal(j)
	if err != nil {
		println("Some error")
	}
	w.Write(res)
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/", a.home).Methods("GET")
	a.Router.HandleFunc("/products", a.products.GetProducts).Methods("POST")
	a.Router.HandleFunc("/product", a.products.CreateProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.products.GetProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.products.UpdateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.products.DeleteProduct).Methods("DELETE")

	a.Router.HandleFunc("/stores/{id:[0-9]+}/products", a.GetProducts).Methods("GET")
	a.Router.HandleFunc("/stores/{id:[0-9]+}", a.AddProducts).Methods("POST")
	a.Router.HandleFunc("/stores/buyProduct", a.BuyProduct).Methods("POST")

	a.Router.HandleFunc("/recommendation/store/getTopProducts/{storeId:[0-9]+}", a.TopProductsInStore).Methods("GET")
	a.Router.HandleFunc("/recommendation/topProductsOFAllStores", a.TopProductsForAllStores).Methods("GET")

}

func (a *App) GetProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	products, err := a.stores.GetProducts(id, limit, start)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Some Error Occurred")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, products)
}

func (a *App) AddProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 64, 10)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}
	var products []models.ProductModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&products); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	payload, err := a.stores.AddProducts(id, products)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Some Error Occurred")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, payload)
}

func (a *App) BuyProduct(w http.ResponseWriter, r *http.Request) {
	productId, _ := strconv.ParseInt(r.FormValue("productId"), 10, 64)
	storeId, _ := strconv.ParseInt(r.FormValue("storeId"), 10, 64)
	if productId == 0 || storeId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store or Product ID")
		return
	}
	payload, err := a.stores.BuyProduct(productId, storeId)
	if err != nil {
		if err.Error() == "record not found" {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, payload)
}

func (a *App) TopProductsInStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeId, err := strconv.ParseInt(vars["storeId"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}
	payload, err := a.orders.TopProductsInStore(storeId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, payload)
}

func (a *App) TopProductsForAllStores(w http.ResponseWriter, r *http.Request) {

	payload, err := a.orders.TopProductsForAllStores()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, payload)
}
