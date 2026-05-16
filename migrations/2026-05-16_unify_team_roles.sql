-- Миграция: унификация ролей в командах и проектах
-- Дата: 2026-05-16

-- 1. Приводим устаревшие роли команд к проектным ролям
UPDATE public.team_members
SET role = 'project_lead'
WHERE role = 'admin';

UPDATE public.team_members
SET role = 'researcher'
WHERE role = 'member';

-- 2. Приводим устаревшие строковые роли участников проектов
UPDATE public.project_members
SET role = 'project_lead'
WHERE role IN ('manager', 'lead');

UPDATE public.project_members
SET role = 'viewer'
WHERE role = 'participant';

-- 3. Синхронизируем role_id для project_members на основе строковой роли
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT id, name FROM public.roles WHERE is_system = false LOOP
        UPDATE public.project_members
        SET role_id = r.id
        WHERE role = r.name AND (role_id IS NULL OR role_id != r.id);
    END LOOP;
END $$;

-- 4. Для project_members без role_id ставим viewer как безопасный fallback
UPDATE public.project_members
SET role = 'viewer', role_id = (SELECT id FROM public.roles WHERE name = 'viewer')
WHERE role_id IS NULL;
