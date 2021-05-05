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
    FOREIGN KEY (shop_id) REFERENCES shops (id) ON DELETE CASCADE
);