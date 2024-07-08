package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHomeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/", homeHandler)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Welcome to XYZ Online Retail!", rr.Body.String())
}

func TestGetProductsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/products", getProductsHandler)

	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var products []Product
	json.NewDecoder(rr.Body).Decode(&products)
	assert.NotEmpty(t, products)
}

func TestAddProductHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/add_product", addProductHandler)

	newProduct := Product{ID: 3, Name: "Tablet", Price: 299.99}
	jsonProduct, _ := json.Marshal(newProduct)
	req, err := http.NewRequest("POST", "/add_product", bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var product Product
	json.NewDecoder(rr.Body).Decode(&product)
	assert.Equal(t, newProduct.Name, product.Name)
}
