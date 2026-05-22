-- Миграция для исправления таблицы task_history
-- Дата: 2026-05-13 (обновлено)

-- Удаляем устаревшие колонки
ALTER TABLE task_history DROP COLUMN IF EXISTS user_id;
ALTER TABLE task_history DROP COLUMN IF EXISTS old_status;
ALTER TABLE task_history DROP COLUMN IF EXISTS new_status;

-- Убедимся, что нужные колонки существуют
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS changed_by INTEGER;
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS field_name VARCHAR(100);
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS old_value TEXT;
ALTER TABLE task_history ADD COLUMN IF NOT EXISTS new_value TEXT;

SELECT 'Миграция task_history выполнена' as result;
