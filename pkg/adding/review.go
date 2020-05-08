package adding

import (
	"time"
)

// Review model
type Review struct {
	User      User      `json:"user"`
	Stars     int       `json:"stars" `
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
