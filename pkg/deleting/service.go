/*
Package deleting includes database deleting operations
*/
package deleting

import (
	"github.com/jinzhu/gorm"
)

// Service provides models deleting operations.
type Service interface {
	Delete(interface{}, string, *gorm.DB)
}

// Delete takes an item of the specified model from the database and permanently deletes it
func Delete(model interface{}, id string, db *gorm.DB) error {
	if err := db.Delete(model, id).Error; err != nil {
		return err
	}
	return nil
}
