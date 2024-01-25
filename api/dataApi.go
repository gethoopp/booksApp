package api

import (
	"github.com/gethoopp/booksApp/services"
	"github.com/gin-gonic/gin"
)

func RestGet() {
	r := gin.Default()
	res := services.Getdata
	req := services.UpdateHandler
	r.GET("/get", res)

	r.PUT("/update/:id", req)

	r.Run(":8080")
}
