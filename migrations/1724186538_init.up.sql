-- Migration Up: init

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name      text NOT NULL,
	email     text NOT NULL,
	phone     text NOT NULL,
	picture   text NOT NULL,
	created_at timestamp NOT NULL
);
CREATE INDEX users_email_trgm_idx ON users USING GIN (email gin_trgm_ops);
CREATE INDEX users_created_at_idx ON users (created_at);