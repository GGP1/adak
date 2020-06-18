/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/GGP1/palo/internal/utils/database"
)

// Service provides models listing operations.
type Service interface {
	GetAll(interface{})
	GetOne(interface{}, string)
}

// GetAll lists all the items of the specified models in the database
func GetAll(model interface{}) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	if err = db.Find(model).Error; err != nil {
		return err
	}
	return nil
}

// GetOne lists just one item of the specified model from the database
func GetOne(model interface{}, id string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	if err := db.First(model, id).Error; err != nil {
		return err
	}
	return nil
}
