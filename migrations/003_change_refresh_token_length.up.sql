ALTER TABLE users RENAME refresh_token TO refresh_token_hash;
ALTER TABLE users ALTER COLUMN refresh_token_hash TYPE VARCHAR(60);
