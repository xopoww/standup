ALTER TABLE users ADD COLUMN whitelisted BOOLEAN DEFAULT false;

UPDATE users SET whitelisted = uw.whitelisted FROM (SELECT user_id, whitelisted FROM user_whitelist) as uw WHERE id = uw.user_id;

DROP TABLE user_whitelist;