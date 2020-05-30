/*
Package rest contains all the functions related to the rest api
*/
package rest

import (
	"fmt"

	h "github.com/GGP1/palo/pkg/http/rest/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter returns gin Engine
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		cookie, err := c.Request.Cookie("sessionID")
		if err != nil {
			c.SetCookie("WelcomeCookie", "1", 0, "/", "localhost", false, true)
		}
		fmt.Println(cookie)
		c.String(200, "Welcome to my golang backend server!")
	})

	router.POST("/login", h.Login)

	product := router.Group("/products")
	{
		product.GET("/", h.GetProducts)
		product.POST("/add", h.AddProduct)
		product.GET("/:id", h.GetAProduct)
		product.PUT("/:id", h.UpdateProduct)
		product.DELETE("/:id", h.DeleteProduct)
	}

	user := router.Group("/users")
	{
		user.GET("/", h.GetUsers)
		user.POST("/add", h.AddUser)
		user.GET("/:id", h.GetAUser)
		user.PUT("/:id", h.UpdateUser)
		user.DELETE("/:id", h.DeleteUser)
	}

	review := router.Group("/reviews")
	{
		review.GET("/", h.GetReviews)
		review.POST("/add", h.AddReview)
		review.GET("/:id", h.GetAReview)
		review.DELETE("/:id", h.DeleteReview)
	}

	shop := router.Group("/shops")
	{
		shop.GET("/", h.GetShops)
		shop.POST("/add", h.AddShop)
		shop.GET("/:id", h.GetAShop)
		shop.PUT("/:id", h.UpdateShop)
		shop.DELETE("/:id", h.DeleteShop)
	}

	return router
}
