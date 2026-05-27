-- Таблица lineage: кто из каких датасетов породил новый датасет через задачу
CREATE TABLE IF NOT EXISTS public.dataset_dependencies
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    source_dataset_id integer NOT NULL,
    target_dataset_id integer NOT NULL,
    task_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT dataset_dependencies_pkey PRIMARY KEY (id),
    CONSTRAINT uq_dd UNIQUE (source_dataset_id, target_dataset_id, task_id),
    CONSTRAINT fk_dd_source FOREIGN KEY (source_dataset_id)
        REFERENCES public.datasets (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_dd_target FOREIGN KEY (target_dataset_id)
        REFERENCES public.datasets (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_dd_task FOREIGN KEY (task_id)
        REFERENCES public.tasks (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_dd_source ON public.dataset_dependencies(source_dataset_id);
CREATE INDEX IF NOT EXISTS idx_dd_target ON public.dataset_dependencies(target_dataset_id);
CREATE INDEX IF NOT EXISTS idx_dd_task ON public.dataset_dependencies(task_id);

-- Добавляем created_by в datasets (для аудита)
ALTER TABLE public.datasets ADD COLUMN IF NOT EXISTS created_by integer;

ALTER TABLE IF EXISTS public.datasets
    ADD CONSTRAINT fk_dataset_created_by FOREIGN KEY (created_by)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;
