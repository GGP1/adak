/*
Package handlers contains all the functions used by the api
*/
package handlers

import (
	"net/http"

	"palo/pkg/adding"
	"palo/pkg/deleting"
	"palo/pkg/listing"
	"palo/pkg/model"
	"palo/pkg/updating"

	"github.com/gin-gonic/gin"
)

// GetProducts func
func GetProducts(c *gin.Context) {
	var product []model.Product
	err := listing.GetProducts(&product)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, product)
}

// GetAProduct func
func GetAProduct(c *gin.Context) {
	var product model.Product
	id := c.Params.ByName("id")

	err := listing.GetAProduct(&product, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, product)
}

// AddProduct func
func AddProduct(c *gin.Context) {
	var product model.Product
	c.BindJSON(&product)

	err := adding.AddProduct(&product)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct func
func UpdateProduct(c *gin.Context) {
	var product model.Product
	id := c.Params.ByName("id")

	err := updating.UpdateProduct(&product, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct func
func DeleteProduct(c *gin.Context) {
	var product model.Product
	id := c.Params.ByName("id")

	err := deleting.DeleteProduct(&product, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, product)
}
