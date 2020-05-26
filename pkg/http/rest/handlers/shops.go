package handlers

import (
	"net/http"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/updating"

	"github.com/gin-gonic/gin"
)

// GetShops lists all the shops
func GetShops(c *gin.Context) {
	var shop []model.Shop
	err := listing.GetShops(&shop)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, shop)
}

// GetAShop lists one shop based on the id
func GetAShop(c *gin.Context) {
	var shop model.Shop
	id := c.Params.ByName("id")

	err := listing.GetAShop(&shop, id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	if shop.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Shop not found"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// AddShop creates a new shop and saves it
func AddShop(c *gin.Context) {
	var shop model.Shop
	c.BindJSON(&shop)

	err := adding.AddShop(&shop)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusCreated, shop)
}

// UpdateShop updates a shop
func UpdateShop(c *gin.Context) {
	var shop model.Shop
	id := c.Params.ByName("id")

	err := updating.UpdateShop(&shop, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shop updated"})
}

// DeleteShop deletes a shop
func DeleteShop(c *gin.Context) {
	var shop model.Shop
	id := c.Params.ByName("id")

	err := deleting.DeleteShop(&shop, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shop deleted"})
}
