-- Миграция: добавление колонки doi для хранения DOI публикаций
-- Дата: 2026-05-16

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS doi TEXT;
