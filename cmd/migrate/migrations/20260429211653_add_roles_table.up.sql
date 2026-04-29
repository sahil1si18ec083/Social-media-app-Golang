BEGIN;

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO roles (id, name, description)
VALUES
    (1, 'creator', 'User can create posts'),
    (2, 'editor', 'User can create and update posts'),
    (3, 'admin', 'User can create, update, and delete posts')
ON CONFLICT (id) DO NOTHING;

ALTER TABLE users
ADD COLUMN role_id BIGINT;

UPDATE users
SET role_id = 1
WHERE role_id IS NULL;

ALTER TABLE users
ALTER COLUMN role_id SET DEFAULT 1,
ALTER COLUMN role_id SET NOT NULL;

ALTER TABLE users
ADD CONSTRAINT users_role_id_fkey
FOREIGN KEY (role_id) REFERENCES roles(id);

COMMIT;
