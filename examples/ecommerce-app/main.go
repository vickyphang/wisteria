package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999.99},
	{ID: 2, Name: "Smartphone", Price: 499.99},
}

func main() {
	router := gin.Default()

	router.GET("/", homeHandler)
	router.GET("/products", getProductsHandler)
	router.GET("/product/:id", getProductHandler)
	router.POST("/add_product", addProductHandler)

	router.Run(":8080")
}

func homeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to XYZ Online Retail!")
}

func getProductsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func getProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	for _, product := range products {
		if product.ID == id {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

func addProductHandler(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products = append(products, newProduct)
	c.JSON(http.StatusCreated, newProduct)
}
