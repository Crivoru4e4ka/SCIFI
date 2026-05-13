-- Миграция для новых таблиц и полей, соответствующих обновленной схеме
-- Дата: 2026-05-10

CREATE TABLE IF NOT EXISTS datasets (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50),
    data_url TEXT,
    parameters JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS experiment_datasets (
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    dataset_id INTEGER NOT NULL REFERENCES datasets(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, dataset_id)
);

CREATE TABLE IF NOT EXISTS hypotheses (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'proposed',
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS task_tags (
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);

-- Новые поля для tasks
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS type TEXT;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS hypothesis_id INTEGER REFERENCES hypotheses(id) ON DELETE SET NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS resource_id INTEGER;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS conclusion TEXT DEFAULT '';
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS parameters JSONB;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS results JSONB;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS artifact_url TEXT;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS metrics JSONB;

-- Кастомные статусы для научного workflow:
-- Draft, Running, Analysis, Peer Review, Archived, Reported
