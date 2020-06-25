/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/GGP1/palo/internal/cfg"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Service provides models listing operations.
type Service interface {
	GetAll(interface{})
	GetOne(interface{}, string)
}

// GetAll lists all the items of the specified models in the database
func GetAll(model interface{}) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Find(model).Error; err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}

// GetOne lists just one item of the specified model from the database
func GetOne(model interface{}, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.First(model, id).Error; err != nil {
		return errors.Wrap(err, "error")
	}
	return nil
}
