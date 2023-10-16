CREATE TABLE IF NOT EXISTS watches (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), title text NOT NULL,
    year integer NOT NULL,
    Price integer NOT NULL,
    watchesType text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
    );