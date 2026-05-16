-- Миграция: RBAC-система для научно-исследовательской платформы
-- Дата: 2026-05-16

-- 1. Роли (системные и проектные)
CREATE TABLE IF NOT EXISTS public.roles (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    is_system boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT roles_pkey PRIMARY KEY (id),
    CONSTRAINT roles_name_key UNIQUE (name)
);

-- 2. Права (permissions)
CREATE TABLE IF NOT EXISTS public.permissions (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    code character varying(100) COLLATE pg_catalog."default" NOT NULL,
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    category character varying(50) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT permissions_pkey PRIMARY KEY (id),
    CONSTRAINT permissions_code_key UNIQUE (code)
);

-- 3. Связь ролей и прав
CREATE TABLE IF NOT EXISTS public.role_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_rp_role FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_rp_permission FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE
);

-- 4. Добавляем role_id в project_members
ALTER TABLE public.project_members ADD COLUMN IF NOT EXISTS role_id integer;
ALTER TABLE public.project_members
    ADD CONSTRAINT fk_pm_role FOREIGN KEY (role_id)
    REFERENCES public.roles (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;

-- 5. Audit log
CREATE TABLE IF NOT EXISTS public.audit_log (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    user_id integer,
    project_id integer,
    action character varying(100) COLLATE pg_catalog."default" NOT NULL,
    entity_type character varying(50) COLLATE pg_catalog."default",
    entity_id integer,
    details text COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT audit_log_pkey PRIMARY KEY (id)
);

ALTER TABLE public.audit_log
    ADD CONSTRAINT fk_audit_user FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;

ALTER TABLE public.audit_log
    ADD CONSTRAINT fk_audit_project FOREIGN KEY (project_id)
    REFERENCES public.projects (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_audit_log_project ON public.audit_log(project_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_user ON public.audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created ON public.audit_log(created_at);

-- 6. Seed: системные роли
INSERT INTO public.roles (name, description, is_system) VALUES
    ('admin', 'Системный администратор', true),
    ('user', 'Обычный пользователь', true),
    ('guest', 'Гость', true)
ON CONFLICT (name) DO NOTHING;

-- 7. Seed: проектные роли
INSERT INTO public.roles (name, description, is_system) VALUES
    ('project_lead', 'Руководитель проекта', false),
    ('scientific_supervisor', 'Научный руководитель', false),
    ('researcher', 'Исследователь', false),
    ('analyst', 'Аналитик', false),
    ('developer', 'Разработчик', false),
    ('reviewer', 'Рецензент', false),
    ('viewer', 'Наблюдатель', false)
ON CONFLICT (name) DO NOTHING;

-- 8. Seed: права (permissions)
INSERT INTO public.permissions (code, name, category) VALUES
    -- Проект
    ('project.view', 'Просмотр проекта', 'project'),
    ('project.edit', 'Редактирование проекта', 'project'),
    ('project.delete', 'Удаление проекта', 'project'),
    ('project.manage_members', 'Управление участниками проекта', 'project'),
    -- Задачи
    ('task.create', 'Создание задач', 'task'),
    ('task.edit', 'Редактирование задач', 'task'),
    ('task.delete', 'Удаление задач', 'task'),
    ('task.assign', 'Назначение исполнителей', 'task'),
    ('task.change_status', 'Изменение статуса задач', 'task'),
    -- Исследования
    ('hypothesis.edit', 'Редактирование гипотез', 'research'),
    ('experiment.create', 'Создание экспериментов', 'research'),
    ('experiment.approve', 'Утверждение экспериментов', 'research'),
    ('analysis.perform', 'Проведение анализа', 'research'),
    ('results.validate', 'Валидация результатов', 'research'),
    -- Данные
    ('dataset.upload', 'Загрузка датасетов', 'data'),
    ('dataset.export', 'Экспорт датасетов', 'data'),
    -- Публикации
    ('publication.create', 'Создание публикаций', 'publication'),
    ('publication.review', 'Рецензирование публикаций', 'publication'),
    ('publication.approve', 'Утверждение публикаций', 'publication'),
    -- Аудит
    ('audit.view', 'Просмотр аудит-лога', 'audit')
ON CONFLICT (code) DO NOTHING;

-- 9. Назначаем права ролям

-- admin: все права
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- user: базовые права
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'user' AND p.code IN ('project.view', 'project.edit', 'task.create', 'task.edit', 'task.change_status', 'hypothesis.edit', 'experiment.create', 'analysis.perform', 'dataset.upload', 'publication.create')
ON CONFLICT DO NOTHING;

-- guest: просмотр + рецензирование
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'guest' AND p.code IN ('project.view', 'publication.review')
ON CONFLICT DO NOTHING;

-- project_lead: все проектные права
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'project_lead'
ON CONFLICT DO NOTHING;

-- scientific_supervisor: утверждение + управление
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'scientific_supervisor' AND p.code IN (
    'project.view', 'project.edit', 'project.manage_members',
    'task.create', 'task.edit', 'task.assign', 'task.change_status',
    'hypothesis.edit', 'experiment.create', 'experiment.approve', 'analysis.perform', 'results.validate',
    'dataset.upload', 'dataset.export',
    'publication.create', 'publication.review', 'publication.approve',
    'audit.view'
)
ON CONFLICT DO NOTHING;

-- researcher: исследовательские задачи
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'researcher' AND p.code IN (
    'project.view', 'task.create', 'task.edit', 'task.change_status',
    'hypothesis.edit', 'experiment.create', 'analysis.perform',
    'dataset.upload', 'publication.create'
)
ON CONFLICT DO NOTHING;

-- analyst: аналитика + визуализации
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'analyst' AND p.code IN (
    'project.view', 'task.create', 'task.edit', 'task.change_status',
    'analysis.perform', 'results.validate', 'dataset.upload', 'dataset.export'
)
ON CONFLICT DO NOTHING;

-- developer: разработка
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'developer' AND p.code IN (
    'project.view', 'task.create', 'task.edit', 'task.change_status',
    'dataset.upload'
)
ON CONFLICT DO NOTHING;

-- reviewer: просмотр + рецензирование
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'reviewer' AND p.code IN (
    'project.view', 'publication.review', 'task.edit'
)
ON CONFLICT DO NOTHING;

-- viewer: только просмотр
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.name = 'viewer' AND p.code = 'project.view'
ON CONFLICT DO NOTHING;
