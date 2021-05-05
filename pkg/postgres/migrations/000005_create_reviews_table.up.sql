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