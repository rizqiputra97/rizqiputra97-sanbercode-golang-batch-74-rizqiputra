package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Struct Bioskop
type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama" binding:"required"`
	Lokasi string  `json:"lokasi" binding:"required"`
	Rating float64 `json:"rating" binding:"required"`
}

var db *sql.DB

func main() {
	var err error

	// Konfigurasi koneksi database PostgreSQL
	connStr := "host=localhost port=5432 user=postgres password=Jatinangor97 dbname=bioskopdb; sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak bisa diakses:", err)
	}

	log.Println("Berhasil terkoneksi ke database!")

	// Inisialisasi Gin
	r := gin.Default()

	// Endpoint
	r.POST("/bioskop", createBioskop)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

// Handler untuk menambahkan bioskop
func createBioskop(c *gin.Context) {
	var input Bioskop

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format JSON tidak valid",
		})
		return
	}

	// Validasi lokasi hanya Jatos atau Jatinangor
	if input.Lokasi != "Jatos" && input.Lokasi != "Jatinangor" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lokasi hanya boleh Jatos atau Jatinangor",
		})
		return
	}

	// Insert ke database
	query := `INSERT INTO bioskop (nama, lokasi, rating) 
	          VALUES ($1, $2, $3) RETURNING id`

	err := db.QueryRow(query, input.Nama, input.Lokasi, input.Rating).Scan(&input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan data",
		})
		return
	}

	c.JSON(http.StatusCreated, input)
}