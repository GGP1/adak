/*
Package rest contains all the functions related to the restful api
*/
package rest

import (
	"fmt"
	h "palo/pkg/http/rest/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SetupRouter returns single and group of routes with their handlers
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		cookie, err := c.Request.Cookie("sessionID")
		if err != nil {
			id := uuid.New()
			c.SetCookie("sessionID", id.String(), 0, "/", "localhost", false, true)
		}
		fmt.Println(cookie)
		c.String(200, "Welcome to my golang backend server!")
	})

	p := router.Group("/products")
	{
		p.GET("/", h.GetProducts)
		p.POST("/add", h.AddProduct)
		p.GET("/:id", h.GetAProduct)
		p.PUT("/:id", h.UpdateProduct)
		p.DELETE("/:id", h.DeleteProduct)
	}

	u := router.Group("/users")
	{
		u.GET("/", h.GetUsers)
		u.POST("/add", h.AddUser)
		u.GET("/:id", h.GetAUser)
		u.PUT("/:id", h.UpdateUser)
		u.DELETE("/:id", h.DeleteUser)
	}

	rw := router.Group("/reviews")
	{
		rw.GET("/", h.GetReviews)
		rw.POST("/add", h.AddReview)
		rw.GET("/:id", h.GetAReview)
		rw.DELETE("/:id", h.DeleteReview)
	}

	return router
}
