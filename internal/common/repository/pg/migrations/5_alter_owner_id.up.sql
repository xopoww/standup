BEGIN;

ALTER TABLE messages RENAME COLUMN owner_id TO owner_name;
ALTER TABLE messages ADD COLUMN owner_id bigint;

UPDATE messages SET owner_id = u.id FROM (SELECT id, username FROM users) as u WHERE owner_name = u.username;
DELETE FROM messages WHERE owner_id IS NULL;
ALTER TABLE messages ALTER COLUMN owner_id SET NOT NULL;
ALTER TABLE messages DROP COLUMN owner_name;

-- Also fix constraints for users
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD PRIMARY KEY (id);

COMMIT;