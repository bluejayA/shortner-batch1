-- URL 단축 서비스 초기 스키마

CREATE TABLE IF NOT EXISTS urls (
    slug        VARCHAR(16)  PRIMARY KEY,
    original    TEXT         NOT NULL,
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id          SERIAL       PRIMARY KEY,
    key_hash    VARCHAR(64)  NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    revoked_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS click_stats (
    slug        VARCHAR(16)  PRIMARY KEY REFERENCES urls(slug) ON DELETE CASCADE,
    click_count BIGINT       DEFAULT 0,
    updated_at  TIMESTAMPTZ  DEFAULT NOW()
);
