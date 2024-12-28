package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "host=k3n9EG8zPxrgH6Lbzb3gznCZs2R4QMh5@dpg-ctnaut5umphs73c4p70g-a.singapore-postgres.render.com user=final_user password=k3n9EG8zPxrgH6Lbzb3gznCZs2R4QMh5 dbname=final_db_1jle port=5432 sslmode=require")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.POST("/products", func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.QueryRow(
			"INSERT INTO products (name, description, price, image_url) VALUES ($1, $2, $3, $4) RETURNING id",
			product.Name, product.Description, product.Price, product.ImageURL,
		).Scan(&product.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	})

	r.GET("/products", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, description, price, image_url FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var product Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			products = append(products, product)
		}
		c.JSON(http.StatusOK, products)
	})

	r.Run(":8084")
}
