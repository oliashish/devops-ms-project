package main

import (
	"log"
	"net/http"
	"strconv"
	"sync" // For thread-safe map access and ID generation

	"github.com/gin-gonic/gin"
)

// Product represents a product entity
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// In-memory storage for products
var (
	products                 = make(map[string]Product)
	nextProductID int64      = 1
	productsMutex sync.Mutex // Mutex to protect access to 'products' map and 'nextProductID'
)

// init function to populate some dummy data when the service starts
func init() {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	addProduct(Product{Name: "Laptop", Description: "Powerful computing device", Price: 1200.00})
	addProduct(Product{Name: "Mouse", Description: "Ergonomic wireless mouse", Price: 25.50})
	addProduct(Product{Name: "Keyboard", Description: "Mechanical gaming keyboard", Price: 75.00})
}

// addProduct is a helper to add a product with an auto-generated ID
func addProduct(product Product) {
	product.ID = strconv.FormatInt(nextProductID, 10)
	products[product.ID] = product
	nextProductID++
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default() // Creates a Gin router with default middleware (logger and recovery)

	// Define API routes
	router.GET("/products", getProducts)        // Get all products
	router.GET("/products/:id", getProductByID) // Get a product by ID
	router.POST("/products", createProduct)     // Create a new product

	// Start the server on port 8081
	log.Println("Product Service starting on port 8081...")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Product Service failed to start: %v", err)
	}
}

func getProducts(c *gin.Context) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	// Convert map to slice for JSON response
	var productList []Product
	for _, product := range products {
		productList = append(productList, product)
	}
	c.JSON(http.StatusOK, productList)
}

func getProductByID(c *gin.Context) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	id := c.Param("id") // Get the product ID from the URL parameter
	product, exists := products[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	var newProduct Product
	// Bind JSON request body to newProduct struct
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate ID and add to map
	addProduct(newProduct)
	c.JSON(http.StatusCreated, newProduct)
}
