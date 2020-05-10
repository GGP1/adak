package rest

import (
	"net/http"

	"palo/pkg/adding"
	"palo/pkg/deleting"
	"palo/pkg/listing"
	"palo/pkg/models"
	"palo/pkg/updating"

	"github.com/gin-gonic/gin"
)

// GetUsers func
func GetUsers(c *gin.Context) {
	var user []models.User
	err := listing.GetUsers(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// GetAUser func
func GetAUser(c *gin.Context) {
	var user models.User
	id := c.Params.ByName("id")

	err := listing.GetAUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// AddUser func
func AddUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)

	err := adding.AddUser(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser func
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Params.ByName("id")

	err := updating.UpdateUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser func
func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Params.ByName("id")

	err := deleting.DeleteUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}
