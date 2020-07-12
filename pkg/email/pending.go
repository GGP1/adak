package email

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// PendingList contains users that need to validate their email
type PendingList struct {
	DB *gorm.DB
	sync.RWMutex

	Email string `json:"email"`
	Token string `json:"token"`
	r     Repository
}

// NewPendingList creates the email pending list service
func NewPendingList(db *gorm.DB, r Repository) Service {
	return &PendingList{
		DB:    db,
		Email: "",
		Token: "",
		r:     r,
	}
}

// Add pending user to the list
func (p *PendingList) Add(email, token string) error {

	p.Lock()
	p.Email = email
	p.Token = token
	p.Unlock()

	if err := p.DB.Create(p).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the pending list")
	}
	return nil
}

// Read pending emails list
func (p *PendingList) Read() (map[string]string, error) {
	if err := p.DB.Find(p).Error; err != nil {
		return nil, errors.Wrap(err, "error: pending list not found")
	}
	maps := make(map[string]string)
	maps[p.Email] = p.Token

	return maps, nil
}

// Remove deletes a key from the map
func (p *PendingList) Remove(key string) error {
	if err := p.DB.Delete(p, key).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the email from the list")
	}
	return nil
}
