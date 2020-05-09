package rest

import (
	"net/http"

	"palo/pkg/adding"
	"palo/pkg/deleting"
	"palo/pkg/listing"

	"github.com/gin-gonic/gin"
)

// GetReviews func
func GetReviews(c *gin.Context) {
	var review []listing.Review
	err := listing.GetReviews(&review)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}

// GetAReview func
func GetAReview(c *gin.Context) {
	var review listing.Review
	id := c.Params.ByName("id")

	err := listing.GetAReview(&review, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}

// AddReview func
func AddReview(c *gin.Context) {
	var review adding.Review
	c.BindJSON(&review)

	err := adding.AddReview(&review)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}

// DeleteReview func
func DeleteReview(c *gin.Context) {
	var review deleting.Review
	id := c.Params.ByName("id")

	err := deleting.DeleteReview(&review, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}
