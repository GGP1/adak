/*
Package deleting includes database deleting operations
*/
package deleting

import (
	"github.com/GGP1/palo/internal/cfg"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Service provides models deleting operations.
type Service interface {
	Delete(interface{}, string)
}

// Delete takes an item of the specified model from the database and permanently deletes it
func Delete(model interface{}, id string) error {
	db, err := gorm.Open("postgres", cfg.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Delete(model, id).Error; err != nil {
		return errors.Wrapf(err, "error: couldn't delete the %s", model)
	}
	return nil
}
