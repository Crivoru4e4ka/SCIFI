-- Миграция: исправление CHECK constraint для системных ролей пользователей
-- Дата: 2026-05-16

-- Удаляем старый constraint, если он есть (может разрешать другие значения)
ALTER TABLE IF EXISTS public.users
    DROP CONSTRAINT IF EXISTS users_role_check;

-- Добавляем constraint с правильными системными ролями:
-- admin  — системный администратор
-- user   — обычный сотрудник/студент
-- guest  — внешний эксперт/рецензент
-- Проектные роли хранятся отдельно в project_members.role_id -> roles
ALTER TABLE IF EXISTS public.users
    ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'user', 'guest'));
