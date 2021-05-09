CREATE TABLE IF NOT EXISTS hits
(
    id text NOT NULL,
    footprint text,
    path text,
    url text,
    language text,
    user_agent text,
    referer text,
    date timestamp with time zone,
    CONSTRAINT hits_pkey PRIMARY KEY (id)
);