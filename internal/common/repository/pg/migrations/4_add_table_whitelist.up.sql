CREATE TABLE user_whitelist (
    user_id bigint PRIMARY KEY,
    whitelisted BOOLEAN
);

INSERT INTO user_whitelist (user_id, whitelisted) SELECT id, whitelisted FROM users;

ALTER TABLE users DROP COLUMN whitelisted;