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

// GetUsers func
func GetUsers(c *gin.Context) {
	var user []model.User
	err := listing.GetUsers(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// GetAUser func
func GetAUser(c *gin.Context) {
	var user model.User
	id := c.Params.ByName("id")

	err := listing.GetAUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// AddUser func
func AddUser(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)

	err := adding.AddUser(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser func
func UpdateUser(c *gin.Context) {
	var user model.User
	id := c.Params.ByName("id")

	err := updating.UpdateUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser func
func DeleteUser(c *gin.Context) {
	var user model.User
	id := c.Params.ByName("id")

	err := deleting.DeleteUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}
