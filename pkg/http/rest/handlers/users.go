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

// GetUsers lists all the users
func GetUsers(c *gin.Context) {
	var user []model.User
	err := listing.GetUsers(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, user)
}

// GetAUser lists one user based on the id
func GetAUser(c *gin.Context) {
	var user model.User
	id := c.Params.ByName("id")

	err := listing.GetAUser(&user, id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// AddUser creates a new user and saves it
func AddUser(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)

	err := adding.AddUser(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context) {
	var user model.User
	id := c.Params.ByName("id")

	err := updating.UpdateUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	var user model.User
	id := c.Params.ByName("id")

	err := deleting.DeleteUser(&user, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
