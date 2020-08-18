// Package storage implements functions for the manipulation
// of databases and caches.
package storage

import (
	"fmt"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/auth/email"
	"github.com/GGP1/palo/pkg/model"
	"github.com/GGP1/palo/pkg/shopping"
	"github.com/GGP1/palo/pkg/shopping/ordering"
	"github.com/GGP1/palo/pkg/tracking"

	"github.com/jinzhu/gorm"
)

// PostgresConnect creates a connection with the database using the postgres driver
// and checks the existence of all the tables.
// It returns a pointer to the gorm.DB struct, the close function and an error.
func PostgresConnect() (*gorm.DB, func() error, error) {
	db, err := gorm.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't open the database: %w", err)
	}

	if err := db.DB().Ping(); err != nil {
		return nil, nil, fmt.Errorf("database connection died: %w", err)
	}

	if err := tableExists(db,
		&model.Product{}, &model.User{}, &model.Review{}, &model.Shop{}, &model.Location{},
		&shopping.Cart{}, &shopping.CartProduct{},
		&tracking.Hit{},
		&ordering.Order{}, &ordering.OrderCart{}, &ordering.OrderProduct{}); err != nil {
		return nil, nil, err
	}

	if db.Table("pending_list").HasTable(&email.List{}) != true {
		db.Table("pending_list").CreateTable(&email.List{}).AutoMigrate(&email.List{})
	}

	if db.Table("validated_list").HasTable(&email.List{}) != true {
		db.Table("validated_list").CreateTable(&email.List{}).AutoMigrate(&email.List{})
	}

	if err := deleteOrdersTrigger(db); err != nil {
		return nil, nil, err
	}

	return db, db.Close, nil
}

// Check if a table is already created, if not, create it.
// Plus model automigration.
func tableExists(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		db.AutoMigrate(model)
		if db.HasTable(model) != true {
			err := db.CreateTable(model).Error
			if err != nil {
				return fmt.Errorf("couldn't create %v table: %w", model, err)
			}
		}
	}
	return nil
}

// If they don't exist, create a function that deletes every order that is outdated,
// giving a margin of 2 days, and create a trigger that executes every time we insert
// new orders.
func deleteOrdersTrigger(db *gorm.DB) error {
	function := `
	IF NOT EXISTS CREATE FUNCTION delete_old_orders() RETURNS trigger
    	LANGUAGE plpgsql
		AS $$
	BEGIN
		DELETE FROM orders WHERE delivery_date < NOW() - INTERVAL '2 days';
		RETURN NULL;
	END;
	$$;`

	trigger := `
	IF NOT EXISTS CREATE TRIGGER trigger_delete_old_orders
		AFTER INSERT ON orders
		EXECUTE PROCEDURE delete_old_orders();`

	if err := db.Raw(function).Error; err != nil {
		return fmt.Errorf("database function: %w", err)
	}

	if err := db.Raw(trigger).Error; err != nil {
		return fmt.Errorf("database trigger: %w", err)
	}

	return nil
}
