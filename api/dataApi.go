package api

import (
	"github.com/gethoopp/booksApp/services"
	"github.com/gin-gonic/gin"
)

func RestGet() {
	r := gin.Default()
	res := services.Getdata

	r.GET("/get", res)

	r.Run(":8080")
}
