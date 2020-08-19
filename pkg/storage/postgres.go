// Package storage implements functions for the manipulation
// of databases and caches.
package storage

import (
	"fmt"

	"github.com/GGP1/palo/internal/cfg"

	"github.com/jmoiron/sqlx"
)

// PostgresConnect creates a connection with the database using the postgres driver
// and checks the existence of all the tables.
// It returns a pointer to the sql.DB struct, the close function and an error.
func PostgresConnect() (*sqlx.DB, func() error, error) {
	db, err := sqlx.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't open the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("database connection died: %w", err)
	}

	db.MustExec(tables)

	// if err := deleteOrdersTrigger(db); err != nil {
	// 	return nil, nil, err
	// }

	return db, db.Close, nil
}

// If they don't exist, create a function that deletes every order that is outdated,
// giving a margin of 2 days, and create a trigger that executes every time we insert
// new orders.
func deleteOrdersTrigger(db *sqlx.DB) error {
	function := `
	CREATE FUNCTION IF NOT EXISTS delete_old_orders() RETURNS trigger
    	LANGUAGE plpgsql
		AS $$
	BEGIN
		DELETE FROM orders WHERE delivery_date < NOW() - INTERVAL '2 days';
		RETURN NULL;
	END;
	$$;`

	trigger := `
	CREATE TRIGGER IF NOT EXISTS trigger_delete_old_orders
		AFTER INSERT ON orders
		EXECUTE PROCEDURE delete_old_orders();`

	_, err := db.Exec(function)
	if err != nil {
		return fmt.Errorf("database function: %w", err)
	}

	_, err = db.Exec(trigger)
	if err != nil {
		return fmt.Errorf("database trigger: %w", err)
	}

	return nil
}

// Order matters
var tables = `
CREATE TABLE IF NOT EXISTS users
(
    id text NOT NULL,
	cart_id text NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS shops
(
    id text NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp with time zone,
    CONSTRAINT shops_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS locations
(
    shop_id text NOT NULL,
    country text NOT NULL,
    state text NOT NULL,
    zip_code text NOT NULL,
    city text NOT NULL,
    address text NOT NULL,
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS products
(
    id text NOT NULL,
    shop_id text NOT NULL,
    stock integer,
    brand text NOT NULL,
    category text NOT NULL,
    type text NOT NULL,
    description text,
    weight numeric NOT NULL,
    taxes numeric,
    discount numeric,
    subtotal numeric NOT NULL,
    total numeric NOT NULL,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp with time zone,
    CONSTRAINT products_pkey PRIMARY KEY (id),
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reviews
(
    id text NOT NULL,
    stars integer NOT NULL,
    comment text,
    user_id text NOT NULL,
    product_id text,
    shop_id text,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp with time zone,
    CONSTRAINT reviews_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS carts
(
    id text NOT NULL,
    counter integer,
    weight numeric,
    discount numeric,
    taxes numeric,
    subtotal numeric,
    total numeric,
    CONSTRAINT carts_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS cart_products
(
    id text NOT NULL,
    cart_id text,
    quantity integer,
    brand text,
    category text,
    type text,
    description text,
    weight numeric,
    taxes numeric,
    discount numeric,
    subtotal numeric,
    total numeric,
    CONSTRAINT cart_products_pkey PRIMARY KEY (id),
    FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS hits
(
    id text NOT NULL,
    footprint text,
    path text,
    url text,
    language text,
    user_agent text,
    referer text,
    date timestamp with time zone,
    CONSTRAINT hits_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS orders
(
    id text NOT NULL,
    user_id text,
    currency text,
    address text,
    city text,
    state text,
    zip_code text,
    country text,
    status text,
    ordered_at timestamp with time zone,
    delivery_date timestamp with time zone,
    cart_id text,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS order_products
(
    id text NOT NULL,
    order_id text,
    product_id integer,
    quantity integer,
    brand text,
    category text,
    type text,
    description text,
    weight numeric,
    discount numeric,
    taxes numeric,
    subtotal numeric,
    total numeric,
    CONSTRAINT order_products_pkey PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS order_carts
(
    order_id text NOT NULL,
    counter integer,
    weight numeric,
    discount numeric,
    taxes numeric,
    subtotal numeric,
    total numeric,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pending_list
(
    email text,
    token text
);

CREATE TABLE IF NOT EXISTS validated_list
(
    email text,
    token text
)
`
