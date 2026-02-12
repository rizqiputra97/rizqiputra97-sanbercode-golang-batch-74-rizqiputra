package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"bioskop-app/models"
)

var DB *sql.DB

func SetDB(database *sql.DB) {
	DB = database
}

// ===================== CREATE =====================
func CreateBioskop(c *gin.Context) {
	var input models.Bioskop

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	if input.Lokasi != "Jatos" && input.Lokasi != "Jatinangor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lokasi hanya boleh Jatos atau Jatinangor"})
		return
	}

	query := `INSERT INTO bioskop (nama, lokasi, rating) 
			  VALUES ($1, $2, $3) RETURNING id`

	err := DB.QueryRow(query, input.Nama, input.Lokasi, input.Rating).Scan(&input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// ===================== READ ALL =====================
func GetAllBioskop(c *gin.Context) {
	rows, err := DB.Query("SELECT id, nama, lokasi, rating FROM bioskop")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var bioskopList []models.Bioskop

	for rows.Next() {
		var b models.Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
			return
		}
		bioskopList = append(bioskopList, b)
	}

	c.JSON(http.StatusOK, bioskopList)
}

// ===================== READ BY ID =====================
func GetBioskopByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var b models.Bioskop
	query := "SELECT id, nama, lokasi, rating FROM bioskop WHERE id=$1"

	err := DB.QueryRow(query, id).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan"})
		return
	}

	c.JSON(http.StatusOK, b)
}

// ===================== UPDATE =====================
func UpdateBioskop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input models.Bioskop
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	if input.Lokasi != "Jatos" && input.Lokasi != "Jatinangor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lokasi hanya boleh Jatos atau Jatinangor"})
		return
	}

	query := `UPDATE bioskop 
			  SET nama=$1, lokasi=$2, rating=$3 
			  WHERE id=$4`

	result, err := DB.Exec(query, input.Nama, input.Lokasi, input.Rating, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update data"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diupdate"})
}

// ===================== DELETE =====================
func DeleteBioskop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	result, err := DB.Exec("DELETE FROM bioskop WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}