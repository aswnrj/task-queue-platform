CREATE TABLE IF NOT EXISTS tasks (
    id          VARCHAR(36) PRIMARY KEY,
    type        VARCHAR(100) NOT NULL,
    payload     TEXT NOT NULL DEFAULT '',
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    result      TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_status ON tasks(status);