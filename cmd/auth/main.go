package main

import (
	"gorm.io/gorm"

	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"

	"fmt"
	"log"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Username string
	Password string
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

var jwtKey = []byte("secret_key")

type JWTClaim struct {
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT generates JWT token for a given email
func GenerateJWT(email, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:   "auth-service",
			Subject:  username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// HashPassword hashes a plain password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares plain password and hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidateToken validates the JWT token and returns the email
func ValidateToken(signedToken string) (email string, err error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Email, nil
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{Email: input.Email, Password: hashedPassword, Username: input.Username}
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already registered"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password 2"})
		return
	}

	token, err := GenerateJWT(user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Request does not contain an access token"})
			c.Abort()
			return
		}

		email, err := ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Next()
	}
}

func main() {
	dsn := "host=auth-db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&User{})
	DB = db

	r := gin.Default()

	r.POST("/register", Register)
	r.POST("/login", Login)

	protected := r.Group("/protected")
	protected.Use(AuthMiddleware())
	protected.GET("/profile", func(c *gin.Context) {
		email, _ := c.Get("email")
		c.JSON(200, gin.H{"email": email})
	})

	fmt.Println("Server running at http://localhost:8888")
	r.Run(":8888")

}
