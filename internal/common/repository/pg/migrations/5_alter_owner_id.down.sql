BEGIN;

ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD PRIMARY KEY (username);

ALTER TABLE messages ADD COLUMN owner_name VARCHAR(32);
UPDATE messages SET owner_name = u.username FROM (SELECT id, username FROM users) as u WHERE owner_id = u.id;
DELETE FROM messages WHERE owner_name IS NULL;
ALTER TABLE messages ALTER COLUMN owner_name SET NOT NULL;
ALTER TABLE messages DROP COLUMN owner_id;
ALTER TABLE messages RENAME COLUMN owner_name TO owner_id;

COMMIT;