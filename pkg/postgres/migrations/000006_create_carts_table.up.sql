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