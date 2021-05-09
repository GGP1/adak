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