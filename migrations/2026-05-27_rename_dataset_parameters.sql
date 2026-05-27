-- Переименование parameters → metadata для улучшения UX
ALTER TABLE public.datasets RENAME COLUMN parameters TO metadata;
