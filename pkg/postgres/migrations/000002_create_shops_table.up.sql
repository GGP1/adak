CREATE TABLE IF NOT EXISTS shops
(
    id text NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp with time zone,
    CONSTRAINT shops_pkey PRIMARY KEY (id)
);