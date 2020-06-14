/*
Package listing includes database listing operations
*/
package listing

import (
	"github.com/jinzhu/gorm"
)

// Service provides models listing operations.
type Service interface {
	GetAll(interface{}, *gorm.DB)
	GetOne(interface{}, string, *gorm.DB)
}

// GetAll lists all the items of the specified models in the database
func GetAll(model interface{}, db *gorm.DB) error {
	if err := db.Find(model).Error; err != nil {
		return err
	}
	return nil
}

// GetOne lists just one item of the specified model from the database
func GetOne(model interface{}, id string, db *gorm.DB) error {
	if err := db.First(model, id).Error; err != nil {
		return err
	}
	return nil
}
