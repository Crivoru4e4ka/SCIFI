-- Миграция: связь проектов с командами и тип выполнения
-- Дата: 2026-05-16

-- 1. Добавляем team_id в projects (nullable, FK к teams)
ALTER TABLE IF EXISTS public.projects
    ADD COLUMN IF NOT EXISTS team_id INTEGER;

ALTER TABLE IF EXISTS public.projects
    ADD CONSTRAINT fk_project_team FOREIGN KEY (team_id)
    REFERENCES public.teams (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;

-- 2. Добавляем execution_type в projects
ALTER TABLE IF EXISTS public.projects
    ADD COLUMN IF NOT EXISTS execution_type VARCHAR(20) DEFAULT 'manual';

-- Заполняем существующие записи
UPDATE public.projects SET execution_type = 'manual' WHERE execution_type IS NULL;

-- CHECK constraint на допустимые значения
ALTER TABLE IF EXISTS public.projects
    DROP CONSTRAINT IF EXISTS projects_execution_type_check;

ALTER TABLE IF EXISTS public.projects
    DROP CONSTRAINT IF EXISTS projects_execution_type_check;

ALTER TABLE IF EXISTS public.projects
    ADD CONSTRAINT projects_execution_type_check CHECK (execution_type IN ('team', 'manual'));

-- 3. Синхронизируем project_members.role_id для создателей проектов.
--    Создатель получает роль 'lead' (руководитель проекта) = id из roles.
--    Сначала найдём id роли 'lead'.
DO $$
DECLARE
    lead_role_id INT;
BEGIN
    SELECT id INTO lead_role_id FROM public.roles WHERE name = 'project_lead' LIMIT 1;
    
    IF lead_role_id IS NOT NULL THEN
        -- Обновляем role_id для всех участников с ролью 'manager' (старое значение) или 'lead'
        UPDATE public.project_members
        SET role_id = lead_role_id
        WHERE role IN ('manager', 'lead') AND role_id IS NULL;
        
        -- Также обновляем строковую роль к единому project_lead
        UPDATE public.project_members
        SET role = 'project_lead'
        WHERE role IN ('manager', 'lead');
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_projects_team_id ON public.projects(team_id);
CREATE INDEX IF NOT EXISTS idx_projects_execution_type ON public.projects(execution_type);
