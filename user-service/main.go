package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// User represents a user entity
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory storage for users
var (
	users                 = make(map[string]User)
	nextUserID int64      = 1
	usersMutex sync.Mutex // Mutex to protect access to 'users' map and 'nextUserID'
)

// init function to populate some dummy data when the service starts
func init() {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	addUser(User{Name: "Alice Smith", Email: "alice@example.com"})
	addUser(User{Name: "Bob Johnson", Email: "bob@example.com"})
	addUser(User{Name: "Charlie Brown", Email: "charlie@example.com"})
}

// addUser is a helper to add a user with an auto-generated ID
func addUser(user User) {
	user.ID = strconv.FormatInt(nextUserID, 10)
	users[user.ID] = user
	nextUserID++
}

func main() {
	// Set Gin to release mode for production, or debug mode for development
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default() // Creates a Gin router with default middleware (logger and recovery)

	// API routes
	router.GET("/users", getUsers)        // Get all users
	router.GET("/users/:id", getUserByID) // Get a user by ID
	router.POST("/users", createUser)     // Create a new user

	// Start the server on port 8080
	log.Println("User Service starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("User Service failed to start: %v", err)
	}
}

func getUsers(c *gin.Context) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	// Convert map to slice for JSON response
	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}
	c.JSON(http.StatusOK, userList)
}

func getUserByID(c *gin.Context) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	id := c.Param("id") // Get the user ID from the URL parameter
	user, exists := users[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	var newUser User
	// Bind JSON request body to newUser struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate ID and add to map
	addUser(newUser)
	c.JSON(http.StatusCreated, newUser)
}
