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

	// Set the maximum number of open and idle connections
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Create a context for the database operations
	ctx := context.Background()

	// Execute the SELECT query
	rows, err := db.QueryContext(ctx, "SELECT id, details FROM buku")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution error"})
		return
	}
	defer rows.Close()

	var result []gin.H // Use a slice of gin.H to store JSON objects

	for rows.Next() {
		var id int
		var detail string

		// Scan the values from the current row
		err := rows.Scan(&id, &detail)
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

		// Append the result to the slice
		result = append(result, gin.H{"id": id, "details": details})
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, result)
}
