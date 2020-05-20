package handlers

import (
	"net/http"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/model"

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
		c.String(http.StatusNotFound, "Review not found")
	} else {
		c.JSON(http.StatusOK, review)
	}
}
