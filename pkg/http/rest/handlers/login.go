/*
Package handlers contains the methods used by the router
*/
package handlers

import (
	"net/http"

	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/models"
	"github.com/gin-gonic/gin"
)

// Login takes a user and authenticates it
func Login(c *gin.Context) {
	user := models.User{}

	err := c.BindJSON(&user)

	err = user.Validate("login")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
	}

	token, err := auth.SignIn(user.Email, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid email or password")
		return
	}

	c.JSON(http.StatusOK, token)
}
