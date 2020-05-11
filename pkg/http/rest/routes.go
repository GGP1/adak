/*
Package rest handles objects request and response functions aswell as its routes
*/
package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter returns single and group of routes with their handlers
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to my golang server!")
	})

	p := r.Group("/products")
	{
		p.GET("/", GetProducts)
		p.POST("/add", AddProduct)
		p.GET("/:id", GetAProduct)
		p.PUT("/:id", UpdateProduct)
		p.DELETE("/:id", DeleteProduct)
	}

	u := r.Group("/users")
	{
		u.GET("/", GetUsers)
		u.POST("/add", AddUser)
		u.GET("/:id", GetAUser)
		u.PUT("/:id", UpdateUser)
		u.DELETE("/:id", DeleteUser)
	}

	rw := r.Group("/reviews")
	{
		rw.GET("/", GetReviews)
		rw.POST("/add", AddReview)
		rw.GET("/:id", GetAReview)
		rw.DELETE("/:id", DeleteReview)
	}

	return r
}
