package handlers

import (
	"net/http"

	"palo/pkg/adding"
	"palo/pkg/deleting"
	"palo/pkg/listing"
	"palo/pkg/model"

	"github.com/gin-gonic/gin"
)

// GetReviews func
func GetReviews(c *gin.Context) {
	var review []model.Review
	err := listing.GetReviews(&review)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}

// GetAReview func
func GetAReview(c *gin.Context) {
	var review model.Review
	id := c.Params.ByName("id")

	err := listing.GetAReview(&review, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}

// AddReview func
func AddReview(c *gin.Context) {
	var review model.Review
	c.BindJSON(&review)

	err := adding.AddReview(&review)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusCreated, review)
}

// DeleteReview func
func DeleteReview(c *gin.Context) {
	var review model.Review
	id := c.Params.ByName("id")

	err := deleting.DeleteReview(&review, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}
