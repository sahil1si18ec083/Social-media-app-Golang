CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    title TEXT NOT NULL,
    content TEXT NOT NULL,

    tags TEXT[] NOT NULL DEFAULT '{}',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
