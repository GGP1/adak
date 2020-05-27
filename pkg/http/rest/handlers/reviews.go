package handlers

import (
	"net/http"

	"github.com/GGP1/palo/pkg/adding"
	"github.com/GGP1/palo/pkg/deleting"
	"github.com/GGP1/palo/pkg/listing"
	"github.com/GGP1/palo/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetReviews lists all the reviews
func GetReviews(c *gin.Context) {
	var review []models.Review
	err := listing.GetReviews(&review)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, review)
}

// GetAReview lists a review based on the id
func GetAReview(c *gin.Context) {
	var review models.Review
	id := c.Params.ByName("id")

	err := listing.GetAReview(&review, id)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	if review.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// AddReview creates a new review and saves it
func AddReview(c *gin.Context) {
	var review models.Review
	c.BindJSON(&review)

	err := adding.AddReview(&review)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusCreated, review)
}

// DeleteReview deletes a review
func DeleteReview(c *gin.Context) {
	var review models.Review
	id := c.Params.ByName("id")

	err := deleting.DeleteReview(&review, id)

	if err != nil {
		c.String(http.StatusNotFound, "Review not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted"})
}
