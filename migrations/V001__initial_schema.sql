-- V001: Initial schema for kyudo-dojo-hub
-- All tables use UUID primary keys (TEXT) and timestamps with timezone.

CREATE TABLE IF NOT EXISTS dojos (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    address     TEXT NOT NULL,
    target_lanes INTEGER NOT NULL CHECK (target_lanes > 0),
    open_time   TEXT NOT NULL,   -- HH:mm format
    close_time  TEXT NOT NULL,   -- HH:mm format
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    email       TEXT NOT NULL UNIQUE,
    role        TEXT NOT NULL CHECK (role IN ('practitioner', 'manager', 'admin')),
    dan         TEXT CHECK (dan IN ('shodan', 'nidan', 'sandan', 'yondan', 'godan', 'rokudan', 'nanadan', 'hachidan', 'kudan', 'judan')),
    shogo       TEXT CHECK (shogo IN ('renshi', 'kyoshi', 'hanshi')),
    dojo_id     TEXT REFERENCES dojos(id),
    joined_at   TEXT NOT NULL,   -- Date string
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_dojo_id ON users(dojo_id);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS practices (
    id                  TEXT PRIMARY KEY,
    user_id             TEXT NOT NULL REFERENCES users(id),
    dojo_id             TEXT REFERENCES dojos(id),
    date                TEXT NOT NULL,  -- YYYY-MM-DD format
    hit_rate            INTEGER NOT NULL CHECK (hit_rate >= 0 AND hit_rate <= 100),
    arrow_count         INTEGER NOT NULL CHECK (arrow_count >= 1 AND arrow_count <= 1000),
    notes               TEXT NOT NULL DEFAULT '',
    instructor_comment  TEXT NOT NULL DEFAULT '',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_practices_user_id ON practices(user_id);
CREATE INDEX idx_practices_date ON practices(date DESC);

CREATE TABLE IF NOT EXISTS videos (
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL REFERENCES users(id),
    practice_id TEXT REFERENCES practices(id),
    file_name   TEXT NOT NULL,
    file_size   BIGINT NOT NULL CHECK (file_size > 0),
    duration    DOUBLE PRECISION NOT NULL CHECK (duration > 0),
    mime_type   TEXT NOT NULL CHECK (mime_type IN ('video/mp4', 'video/quicktime', 'video/webm')),
    status      TEXT NOT NULL DEFAULT 'completed' CHECK (status IN ('uploading', 'processing', 'completed', 'failed')),
    url         TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_videos_user_id ON videos(user_id);

CREATE TABLE IF NOT EXISTS analyses (
    id             TEXT PRIMARY KEY,
    video_id       TEXT NOT NULL UNIQUE REFERENCES videos(id),
    user_id        TEXT NOT NULL REFERENCES users(id),
    scores         JSONB NOT NULL,   -- HassetsuScores
    phases         JSONB NOT NULL,   -- []PhaseSegment
    overall_score  INTEGER NOT NULL CHECK (overall_score >= 0 AND overall_score <= 100),
    feedback       TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analyses_user_id ON analyses(user_id);
CREATE INDEX idx_analyses_video_id ON analyses(video_id);

CREATE TABLE IF NOT EXISTS reservations (
    id          TEXT PRIMARY KEY,
    dojo_id     TEXT NOT NULL REFERENCES dojos(id),
    user_id     TEXT NOT NULL REFERENCES users(id),
    lane_number INTEGER NOT NULL CHECK (lane_number >= 1),
    date        TEXT NOT NULL,  -- YYYY-MM-DD format
    start_time  TEXT NOT NULL,  -- HH:mm format
    end_time    TEXT NOT NULL,  -- HH:mm format
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Prevent exact duplicate bookings (same lane, same date, same start time)
    CONSTRAINT uq_reservations_lane_time UNIQUE (dojo_id, lane_number, date, start_time)
);

CREATE INDEX idx_reservations_dojo_date ON reservations(dojo_id, date);

CREATE TABLE IF NOT EXISTS exam_checklists (
    id              TEXT PRIMARY KEY,
    user_id         TEXT NOT NULL REFERENCES users(id),
    target_dan      TEXT NOT NULL CHECK (target_dan IN ('shodan', 'nidan', 'sandan', 'yondan', 'godan', 'rokudan', 'nanadan', 'hachidan', 'kudan', 'judan')),
    items           JSONB NOT NULL,   -- []ExamChecklistItem
    progress_rate   INTEGER NOT NULL DEFAULT 0 CHECK (progress_rate >= 0 AND progress_rate <= 100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exam_checklists_user_id ON exam_checklists(user_id);
