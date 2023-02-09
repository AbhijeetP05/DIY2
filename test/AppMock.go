package test

import (
	"DIY2/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// App main app for the program to run, contains a gorm database instance and a router instance for routing
type App struct {
	DB       *gorm.DB
	Router   *mux.Router
	products services.IProducts
	stores   services.IStores
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
	//p := models.ProductModel{ID: 2}
	//p.GetProduct(a.DB)
	//fmt.Println(p.ID, p.Name, p.Price)

	a.products = services.NewProduct(a.DB)
	a.stores = services.NewStore(a.DB)
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

	a.Router.HandleFunc("/stores/{id:[0-9]+}/products", a.stores.GetProducts).Methods("GET")
	a.Router.HandleFunc("/stores/{id:[0-9]+}", a.stores.AddProducts).Methods("POST")

}
