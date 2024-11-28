ALTER TABLE users RENAME refresh_token_hash TO refresh_token;
ALTER TABLE users ALTER COLUMN refresh_token TYPE VARCHAR(100);
