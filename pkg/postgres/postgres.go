// Package postgres implements functions for the manipulation of databases.
package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Connect creates a connection with the database using the postgres driver
// and checks the existence of all the tables.
//
// It returns a pointer to the sql.DB struct, the close function and an error.
func Connect(ctx context.Context, c config.Postgres) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Name, c.SSLMode)

	db, err := sqlx.ConnectContext(ctx, "postgres", url)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open the database")
	}

	if err := Migrate(ctx, db); err != nil {
		return nil, err
	}

	logger.Infof("Connected to postgres on %s", net.JoinHostPort(c.Host, c.Port))
	return db, nil
}

// Migrate creates database tables, indexes and triggers.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	if err := createTables(ctx, db); err != nil {
		return err
	}

	if err := createIndexes(ctx, db); err != nil {
		return err
	}

	return createTriggers(ctx, db)
}

// createTables creates the database tables. It's implemented in a separate function for
// testing purposes.
func createTables(ctx context.Context, db *sqlx.DB) error {
	if _, err := db.ExecContext(ctx, tables); err != nil {
		return errors.Wrap(err, "couldn't create tables")
	}
	return nil
}

// createIndexes creates database indexes.
func createIndexes(ctx context.Context, db *sqlx.DB) error {
	if _, err := db.ExecContext(ctx, indexes); err != nil {
		return errors.Wrap(err, "couldn't create indexes")
	}
	return nil
}

// createTriggers creates database functions and its triggers.
func createTriggers(ctx context.Context, db *sqlx.DB) error {
	if _, err := db.ExecContext(ctx, triggers); err != nil {
		return errors.Wrap(err, "couldn't create triggers")
	}
	return nil
}

// Order matters
const tables = `
CREATE TABLE IF NOT EXISTS users
(
    id text NOT NULL,
	cart_id text NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    verified_email boolean DEFAULT false,
    is_admin boolean DEFAULT false,
    confirmation_code text,
    search tsvector,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp DEFAULT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS shops
(
    id text NOT NULL,
    name text NOT NULL,
    search tsvector,
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
    search tsvector,
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
    CONSTRAINT reviews_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE,
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE
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
    cart_id text NOT NULL,
    quantity integer NOT NULL,
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
    status integer,
    cart_id text,
    created_at timestamp with time zone DEFAULT NOW(),
    ordered_at timestamp with time zone,
    delivery_date timestamp with time zone,
    CONSTRAINT orders_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
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
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS order_products
(
    order_id text NOT NULL,
    product_id text NOT NULL,
    quantity integer,
    brand text,
    category text,
    type text,
    description text,
    weight integer,
    discount integer,
    taxes integer,
    subtotal integer,
    total integer,
    FOREIGN KEY (order_id) 
        REFERENCES orders (id)
        ON DELETE CASCADE
        DEFERRABLE INITIALLY DEFERRED
);`

const indexes = `
CREATE INDEX ON users USING GIN (search);
CREATE INDEX ON shops USING GIN (search);
CREATE INDEX ON products USING GIN (search);

CREATE INDEX ON users (created_at);
CREATE INDEX ON shops (created_at);
CREATE INDEX ON products (created_at);
CREATE INDEX ON reviews (created_at);
CREATE INDEX ON orders (created_at);`

const triggers = `
CREATE OR REPLACE FUNCTION users_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  new.search := to_tsvector('english', new.username);
  return new;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS users_tsvector_update ON users;

CREATE TRIGGER users_tsvector_update BEFORE INSERT OR UPDATE
    ON users FOR EACH ROW EXECUTE PROCEDURE users_tsvector_trigger();

--
    
CREATE OR REPLACE FUNCTION shops_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  new.search := to_tsvector('english', new.name);
  return new;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS shops_tsvector_update ON shops;

CREATE TRIGGER shops_tsvector_update BEFORE INSERT OR UPDATE
    ON shops FOR EACH ROW EXECUTE PROCEDURE shops_tsvector_trigger();

--

CREATE OR REPLACE FUNCTION products_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  new.search :=
  setweight(to_tsvector('english', new.type), 'A')
  || setweight(to_tsvector('english', new.category), 'B')
  || setweight(to_tsvector('english', new.brand), 'C');
  return new;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS products_tsvector_update ON products;

CREATE TRIGGER products_tsvector_update BEFORE INSERT OR UPDATE
    ON products FOR EACH ROW EXECUTE PROCEDURE products_tsvector_trigger();`
