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
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);