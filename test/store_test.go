package test

import (
	"DIY2/models"
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestGetProductsWithWrongID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//mockStore := mocks.NewMockIStores(mockCtrl)
	//id := 5

	req, _ := http.NewRequest("GET", "/stores/5/products", nil)
	//mockStore.GetProducts(nil, req)
	response := executeRequest(req)

	var m []models.ProductModel
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if len(m) != 0 {
		t.Errorf("Expected response size to be 0 but got %q", len(m))
	}
}

func TestGetProductsSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//mockStore := mocks.NewMockIStores(mockCtrl)
	//id := 1

	req, _ := http.NewRequest("GET", "/stores/1/products", nil)
	//mockStore.GetProducts(nil, req)
	response := executeRequest(req)

	var m []models.ProductModel
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if len(m) == 0 {
		t.Errorf("Expected response size to be greater than 0 but got %q", len(m))
	}
}

func TestAddProductsEmpty(t *testing.T) {
	query := "SELECT * FROM products"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	var m []models.ProductModel
	s, _ := json.Marshal(&m)
	//body := io.newReader()
	req, _ := http.NewRequest("POST", "/stores/1", bytes.NewBuffer(s))
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected
	checkResponseCode(t, http.StatusOK, response.Code)
	if finalCount != initialCount {
		t.Errorf("Unnecessary product added.")
	}
}

func TestAddProductsSuccess(t *testing.T) {

	query := "SELECT * FROM products"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	var m []models.ProductModel
	m = append(m, models.ProductModel{Name: "New Test Product", Price: 11.69})
	s, _ := json.Marshal(&m)
	//body := io.newReader()
	req, _ := http.NewRequest("POST", "/stores/1", bytes.NewBuffer(s))
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected
	checkResponseCode(t, http.StatusOK, response.Code)
	if finalCount != initialCount+1 {
		t.Errorf("product not added added.")
	}
}

func TestBuyProductSuccess(t *testing.T) {
	query := "SELECT * FROM orders"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var storeId, productId int64
	storeId = 1
	productId = 2
	form := url.Values{}
	form.Add("storeId", strconv.FormatInt(storeId, 10))
	form.Add("productId", strconv.FormatInt(productId, 10))
	req, _ := http.NewRequest("POST", "/stores/buyProduct", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//response := executePostRequest("/stores/buyProduct/", form)
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected
	checkResponseCode(t, http.StatusOK, response.Code)
	if finalCount != initialCount+1 {
		t.Errorf("order not added.")
	}
}

func TestBuyProduct_NoFormData(t *testing.T) {
	query := "SELECT * FROM orders"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	req, _ := http.NewRequest("POST", "/stores/buyProduct", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	if finalCount != initialCount {
		t.Errorf("order added while it was not expected to.")
	}
}

func TestBuyProductFail(t *testing.T) {
	query := "SELECT * FROM orders"
	result := a.DB.Exec(query)
	initialCount := result.RowsAffected
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var storeId, productId int64
	storeId = 1
	productId = 1
	form := url.Values{}
	form.Add("storeId", strconv.FormatInt(storeId, 10))
	form.Add("productId", strconv.FormatInt(productId, 10))
	req, _ := http.NewRequest("POST", "/stores/buyProduct", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := executeRequest(req)
	result = a.DB.Exec(query)
	finalCount := result.RowsAffected
	checkResponseCode(t, http.StatusNotFound, response.Code)
	if finalCount != initialCount {
		t.Errorf("non existent orderorder added.")
	}
}
