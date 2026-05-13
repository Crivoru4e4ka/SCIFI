-- Миграция для исправления таблицы task_history
-- Дата: 2026-05-13

-- Если колонка user_id существует, переименуем на changed_by
ALTER TABLE task_history RENAME COLUMN user_id TO changed_by;

-- Если нужны другие колонки, добавим их
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS field_name VARCHAR(100);
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS old_value TEXT;
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS new_value TEXT;
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS old_status VARCHAR(50);
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS new_status VARCHAR(50);

SELECT 'Миграция task_history выполнена' as result;
