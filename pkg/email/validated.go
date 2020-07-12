package email

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// ValidatedList contains users that already verified their email
type ValidatedList struct {
	sync.RWMutex
	DB *gorm.DB

	Email string `json:"email"`
	Token string `json:"token"`
	r     Repository
}

// NewValidatedList creates the email validated list service
func NewValidatedList(db *gorm.DB, r Repository) Service {
	return &ValidatedList{
		DB:    db,
		Email: "",
		Token: "",
		r:     r,
	}
}

// Add validated user to the list
func (v *ValidatedList) Add(email, token string) error {

	v.Lock()
	v.Email = email
	v.Token = token
	v.Unlock()

	if err := v.DB.Create(v).Error; err != nil {
		return errors.Wrap(err, "error: couldn't create the validated list")
	}
	return nil
}

// Read validated emails list
func (v *ValidatedList) Read() (map[string]string, error) {
	if err := v.DB.Find(v).Error; err != nil {
		return nil, errors.Wrap(err, "error: validated list not found")
	}
	maps := make(map[string]string)
	maps[v.Email] = v.Token

	return maps, nil
}

// Remove deletes a key from the map
func (v *ValidatedList) Remove(key string) error {
	if err := v.DB.Delete(v, key).Error; err != nil {
		return errors.Wrap(err, "error: couldn't delete the email from the list")
	}
	return nil
}
