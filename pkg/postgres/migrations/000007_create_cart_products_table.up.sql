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
    FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE
);