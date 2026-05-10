-- Миграция для добавления новых полей в таблицу tasks
-- Дата: 2026-05-05

-- Добавить колонку start_date если не существует
ALTER TABLE tasks 
ADD COLUMN IF NOT EXISTS start_date TIMESTAMP NULL;

-- Добавить колонку team если не существует
ALTER TABLE tasks 
ADD COLUMN IF NOT EXISTS team VARCHAR(255) NULL;

-- Вывести результат
SELECT 'Миграция выполнена успешно' as result;
