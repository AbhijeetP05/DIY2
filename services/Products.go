package services

import (
	"DIY2/models"
	"DIY2/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	productRepo models.IProductRepo
}

//go:generate mockgen -destination=../mocks/mock_product.go -package=mocks go-mux/services IProducts
type IProducts interface {
	GetProduct(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
}

func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		log.Fatal("Invalid Product ID")
	}
	productModel := models.ProductModel{ID: &id}
	result := p.productRepo.GetProduct(&productModel)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		utils.RespondWithError(w, http.StatusNotFound, "Product Not Found")
	} else if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
	} else {
		utils.RespondWithJSON(w, http.StatusOK, productModel)
	}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	productModel := models.ProductModel{}
	products, result := p.productRepo.GetProducts(&productModel, limit, start)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
	} else {
		utils.RespondWithJSON(w, http.StatusOK, products)
	}
}

func (p *Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var product models.ProductModel
	if err := decoder.Decode(&product); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	result := p.productRepo.CreateProduct(&product)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, product)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	// get initial product
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	product := models.ProductModel{ID: &id}
	result := p.productRepo.GetProduct(&product)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// get the updated product
	var newProduct models.ProductModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newProduct); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result = p.productRepo.UpdateProduct(&product, &newProduct)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, product)
}

func (p *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// get initial product
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	product := models.ProductModel{ID: &id}
	result := p.productRepo.GetProduct(&product)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	// delete the product
	result = p.productRepo.DeleteProduct(&product)
	if result.Error != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func NewProduct(repo models.IProductRepo) *Products {
	return &Products{productRepo: repo}
}
