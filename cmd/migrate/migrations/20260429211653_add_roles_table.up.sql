BEGIN;

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    level INT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO
  roles (name, description, level)
VALUES
  (
    'user',
    'A user can create posts and comments',
    1
  );

INSERT INTO
  roles (name, description, level)
VALUES
  (
    'moderator',
    'A moderator can update other users posts',
    2
  );

INSERT INTO
  roles (name, description, level)
VALUES
  (
    'admin',
    'An admin can update and delete other users posts',
    3
  );

ALTER TABLE users
ADD COLUMN role_id BIGINT;

UPDATE users
SET role_id = (
    SELECT id
    FROM roles
    WHERE level = 1
)
WHERE role_id IS NULL;

ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;

ALTER TABLE users
ADD CONSTRAINT users_role_id_fkey
FOREIGN KEY (role_id) REFERENCES roles(id);

COMMIT;
