package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var DB *gorm.DB

func initDatabase() {

	dsn := "host=hello-db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})
	user := User{Name: "Jinzhu", Email: "aulia@gmail.com"}
	db.Create(&user) // pass pointer of data to Create

	DB = db
}

// Handler untuk GET all users
func getAllUsers(c *gin.Context) {
	var users []User

	// Query database untuk mengambil semua user
	if err := DB.Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to get users"})
		return
	}

	// Mengembalikan response JSON
	c.JSON(200, users)
}

func main() {

	initDatabase()

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Club !!!!",
		})
	})

	router.GET("/users", getAllUsers)

	router.Run(":7777")

}
