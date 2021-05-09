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
    ordered_at timestamp with time zone,
    delivery_date timestamp with time zone,
    cart_id text,
    CONSTRAINT orders_pkey PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);