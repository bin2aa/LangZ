-- 004: Posts table
DO $$ BEGIN
    CREATE TYPE post_type AS ENUM ('lost_found', 'give_away', 'alert', 'general', 'service_request');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS posts (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id    UUID REFERENCES groups(id) ON DELETE SET NULL,
    type        post_type NOT NULL DEFAULT 'general',
    title       VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    location    GEOGRAPHY(Point, 4326),
    image_urls  TEXT[],
    is_resolved BOOLEAN DEFAULT FALSE,
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_location ON posts USING GIST (location);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_posts_group_id ON posts (group_id);
CREATE INDEX IF NOT EXISTS idx_posts_type ON posts (type);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC);
