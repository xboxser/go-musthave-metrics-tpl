-- Создание таблицы для хранения значений метрик
CREATE TABLE IF NOT EXISTS metrics (
    id    TEXT PRIMARY KEY,
    type  TEXT NOT NULL CHECK (type IN ('gauge', 'counter')),
    delta BIGINT,
    value DOUBLE PRECISION,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание индекса по типу записи 
CREATE INDEX IF NOT EXISTS idx_metrics_type ON metrics(type);