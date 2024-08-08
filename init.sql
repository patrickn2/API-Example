-- CREATE USER clerkuser WITH PASSWORD 'password';
CREATE DATABASE clerk OWNER postgres;
GRANT ALL PRIVILEGES ON DATABASE clerk TO postgres;
\c clerk;
CREATE SCHEMA random_project;
CREATE TABLE random_project.users (
  id BIGSERIAL PRIMARY KEY,
  name      text NOT NULL,
	email     text NOT NULL,
	phone     text NOT NULL,
	picture   text NOT NULL,
	created_at timestamp NOT NULL
);