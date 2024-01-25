package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Getdata(c *gin.Context) {
	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/books")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Create a context for the database operations
	ctx := context.Background()

	/*rowss, erro := db.ExecContext(ctx, `INSERT INTO buku (id,details) VALUES('8,'{"Nama": "My Second Hope", "Image" : "https://i.ibb.co/VxGmX5Y/Blue-Minimalist-Watercolor-Sky-Artistic-Novel-Book-Cover.png", "Author" : "Danna Stroupe"}')`)
	if erro != nil {
		panic(erro)
	}*/

	rows, err := db.QueryContext(ctx, "SELECT id, details, hororCat FROM buku")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution error"})
		return
	}

	var result []gin.H // Use a slice of gin.H to store JSON objects

	for rows.Next() {
		var id int
		var detail string
		var horor string

		// Scan the values from the current row
		err := rows.Scan(&id, &detail, &horor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data scanning error"})
			return
		}

		// Parse the JSON string into a map
		var details map[string]interface{}

		if err := json.Unmarshal([]byte(detail), &details); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error"})
			return
		}

		var hororCat map[string]interface{}

		if err := json.Unmarshal([]byte(horor), &hororCat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error"})
			return

		}

		// Append the result to the slice
		result = append(result, gin.H{"id": id, "details": details, "hororCat": hororCat})
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, result)

	defer rows.Close()
}

type BooksD struct {
	Author string `json:"Author"`
	Image  string `json:"Image"`
	Nama   string `json:"Nama"`
	Status bool   `json:"Status"`
}

func UpdateHandler(c *gin.Context) {

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/books")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Create a context for the database operations
	ctx := context.Background()

	id := c.Param("id")

	// Parse the JSON details from the request body
	var newDetails BooksD
	if err := c.ShouldBindJSON(&newDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Convert the BookDetails struct to a JSON string
	newDetailsJSON, err := json.Marshal(newDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting details to JSON"})
		return
	}

	query := "UPDATE buku SET details = ? WHERE id = ?"
	result, err := db.ExecContext(ctx, query, string(newDetailsJSON), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No rows updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update successful"})
}
