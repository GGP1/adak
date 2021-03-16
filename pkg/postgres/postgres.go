// Package postgres implements functions for the manipulation of databases.
package postgres

import (
	"context"
	"fmt"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Connect creates a connection with the database using the postgres driver
// and checks the existence of all the tables.
//
// It returns a pointer to the sql.DB struct, the close function and an error.
func Connect(ctx context.Context, c *config.Database) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Name, c.SSLMode)

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open the database")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "database connection died")
	}

	db.MustExecContext(ctx, tables)

	logger.Log.Infof("Postgres listening on %s:%s", c.Host, c.Port)

	return db, nil
}

// Order matters
var tables = `
CREATE TABLE IF NOT EXISTS users
(
    id text NOT NULL,
	cart_id text NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    verified_email boolean,
    confirmation_code text,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp DEFAULT NULL,
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
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS products
(
    id text NOT NULL,
    shop_id text NOT NULL,
    stock integer NOT NULL,
    brand text NOT NULL,
    category text NOT NULL,
    type text NOT NULL,
    description text,
    weight integer NOT NULL,
    taxes integer,
    discount integer,
    subtotal integer NOT NULL,
    total integer NOT NULL,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp with time zone,
    CONSTRAINT products_pkey PRIMARY KEY (id),
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE ON UPDATE CASCADE
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
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS carts
(
    id text NOT NULL,
    counter integer,
    weight integer,
    discount integer,
    taxes integer,
    subtotal integer,
    total integer,
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
    weight integer,
    taxes integer,
    discount integer,
    subtotal integer,
    total integer,
    CONSTRAINT cart_products_pkey PRIMARY KEY (id),
    FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE ON UPDATE CASCADE
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
    CONSTRAINT orders_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS order_carts
(
    order_id text NOT NULL,
    counter integer,
    weight integer,
    discount integer,
    taxes integer,
    subtotal integer,
    total integer,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS order_products
(
    product_id text NOT NULL,
    order_id text NOT NULL,
    quantity integer,
    brand text,
    category text,
    type text,
    description text,
    weight integer,
    discount integer,
    taxes integer,
    subtotal integer,
    total integer
);
`