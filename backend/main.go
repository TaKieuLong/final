package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Connected to the database successfully!")
}

func main() {
	initDB()

	r := gin.Default()
	r.GET("/products", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, description, price, image_url FROM products")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch products"})
			return
		}
		defer rows.Close()

		var products []map[string]interface{}
		for rows.Next() {
			var id int
			var name, description, imageURL string
			var price float64
			if err := rows.Scan(&id, &name, &description, &price, &imageURL); err != nil {
				c.JSON(500, gin.H{"error": "Failed to parse products"})
				return
			}
			products = append(products, gin.H{
				"id":          id,
				"name":        name,
				"description": description,
				"price":       price,
				"image_url":   imageURL,
			})
		}
		c.JSON(200, products)
	})

	r.POST("/products", func(c *gin.Context) {
		var product struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
			ImageURL    string  `json:"image_url"`
		}
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO products (name, description, price, image_url) VALUES ($1, $2, $3, $4)",
			product.Name, product.Description, product.Price, product.ImageURL)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(201, gin.H{"message": "Product created successfully!"})
	})

	r.Run(":8084")
}
