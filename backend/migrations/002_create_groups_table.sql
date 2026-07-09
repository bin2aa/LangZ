-- 002: Groups table
CREATE TABLE IF NOT EXISTS groups (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    center_location GEOGRAPHY(Point, 4326),
    radius_meters   INTEGER NOT NULL DEFAULT 1000,
    created_by      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_groups_center_location ON groups USING GIST (center_location);
CREATE INDEX IF NOT EXISTS idx_groups_created_by ON groups (created_by);
