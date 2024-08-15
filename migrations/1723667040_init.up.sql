-- Migration Up: Init

CREATE DATABASE clerk OWNER postgres;
GRANT ALL PRIVILEGES ON DATABASE clerk TO postgres;
\c clerk;
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name      text NOT NULL,
	email     text NOT NULL,
	phone     text NOT NULL,
	picture   text NOT NULL,
	created_at timestamp NOT NULL
);
CREATE INDEX idx_users_created_at_desc_email ON users (created_at DESC, email);
CREATE INDEX idx_users_created_at_email ON users (created_at ASC, email);