/*
Package handlers contains the methods used by the router
*/
package handlers

import (
	"net/http"

	"github.com/GGP1/palo/pkg/auth"
	"github.com/GGP1/palo/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Login takes a user and authenticates it
func Login(c *gin.Context) {
	user := models.User{}

	c.Request.Cookie("SID")

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

	id := uuid.New()
	c.SetCookie("SID", id.String(), 0, "/", "localhost", false, true)

	c.JSON(http.StatusOK, token)
}

// Logout removes the authentication cookie
func Logout(c *gin.Context) {
	c.SetCookie("SID", "0", -1, "/", "localhost", false, true)

	c.String(http.StatusOK, "You are now logged out")
}
