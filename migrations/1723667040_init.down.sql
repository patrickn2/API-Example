-- Migration Down: Init

DROP INDEX idx_users_created_at_desc_email;
DROP INDEX idx_users_created_at_email;

DROP TABLE public.users;
DROP DATABASE clerk;