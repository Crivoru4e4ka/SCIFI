-- Метаданные датасета — простой текст, не JSON blob
ALTER TABLE public.datasets ALTER COLUMN metadata TYPE text;
