-- Создание таблицы контактов
CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    telegram_contact VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы артефактов
CREATE TABLE IF NOT EXISTS artifacts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    description TEXT,
    project_name VARCHAR(255),
    developer_id INTEGER REFERENCES contacts(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE artifacts ADD CONSTRAINT artifacts_type_check   
CHECK (type IN ('table', 'view', 'procedure', 'function', 'index', 'dataset', 'api', 'file'));

-- Создание таблицы полей артефактов
CREATE TABLE IF NOT EXISTS artifact_fields (
    id SERIAL PRIMARY KEY,
    artifact_id INTEGER REFERENCES artifacts(id) ON DELETE CASCADE,
    field_name VARCHAR(255) NOT NULL,
    data_type VARCHAR(100) NOT NULL,
    description TEXT,
    is_pk BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для улучшения производительности
CREATE INDEX idx_artifacts_developer_id ON artifacts(developer_id);
CREATE INDEX idx_artifacts_project_name ON artifacts(project_name);
CREATE INDEX idx_artifact_fields_artifact_id ON artifact_fields(artifact_id);
CREATE INDEX idx_artifacts_created_at ON artifacts(created_at DESC);
CREATE INDEX idx_contacts_created_at ON contacts(created_at DESC);