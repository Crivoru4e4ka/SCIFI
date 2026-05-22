-- Миграция: удаление мёртвых колонок и cleanup после рефакторинга
-- Дата: 2026-05-22

-- Удаляем program_name из projects
ALTER TABLE projects DROP COLUMN IF EXISTS program_name;

-- Удаляем resource_id из tasks (если остался)
ALTER TABLE tasks DROP COLUMN IF EXISTS resource_id;

-- Удаляем audit_log (если ещё существует)
DROP TABLE IF EXISTS public.audit_log CASCADE;

-- Добавляем relation_type в experiment_datasets (если отсутствует)
ALTER TABLE experiment_datasets ADD COLUMN IF NOT EXISTS relation_type VARCHAR(20) DEFAULT 'input';

SELECT 'Cleanup миграция выполнена' as result;
