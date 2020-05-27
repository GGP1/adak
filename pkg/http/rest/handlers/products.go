package handlers

import (
	"net/http"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/models"
	"github.com/GGP1/palo/pkg/updating"

	"github.com/gin-gonic/gin"
)

// GetProducts lists all the products
func GetProducts(c *gin.Context) {
	var product []models.Product
	err := listing.GetProducts(&product)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, product)
}

// GetAProduct lists one product based on the id
func GetAProduct(c *gin.Context) {
	var product models.Product
	id := c.Params.ByName("id")

	err := listing.GetAProduct(&product, id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// AddProduct creates a new product and saves it
func AddProduct(c *gin.Context) {
	var product models.Product
	c.BindJSON(&product)

	err := adding.AddProduct(&product)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates a product
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Params.ByName("id")

	err := updating.UpdateProduct(&product, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

// DeleteProduct deletes a product
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Params.ByName("id")

	err := deleting.DeleteProduct(&product, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
