-- Migration Down: init

DROP INDEX users_email_trgm_idx;
DROP INDEX users_created_at_idx;

DROP TABLE users;
DROP DATABASE clerk;