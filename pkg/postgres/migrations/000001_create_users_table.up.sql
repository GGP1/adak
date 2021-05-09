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
    created_at timestamp with time zone DEFAULT NOW(),
    updated_at timestamp DEFAULT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);