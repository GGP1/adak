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