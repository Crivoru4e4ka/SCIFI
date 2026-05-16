-- Миграция: добавление таблиц грантов и финансирования проектов
-- Дата: 2026-05-16

CREATE TABLE IF NOT EXISTS public.grants (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    title character varying(255) COLLATE pg_catalog."default" NOT NULL,
    code character varying(100) COLLATE pg_catalog."default",
    funding_organization character varying(255) COLLATE pg_catalog."default" NOT NULL,
    country character varying(100) COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    scientific_direction character varying(255) COLLATE pg_catalog."default",
    grant_type character varying(50) COLLATE pg_catalog."default" DEFAULT 'state'::character varying,
    status character varying(50) COLLATE pg_catalog."default" DEFAULT 'draft'::character varying,
    total_amount numeric(15,2) DEFAULT 0,
    currency character varying(10) COLLATE pg_catalog."default" DEFAULT 'RUB'::character varying,
    start_date date,
    end_date date,
    application_deadline date,
    principal_investigator_id integer,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    CONSTRAINT grants_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.project_grant_funding (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    grant_id integer NOT NULL,
    project_id integer NOT NULL,
    section_id integer,
    allocated_amount numeric(15,2) DEFAULT 0,
    funding_purpose text COLLATE pg_catalog."default",
    funding_start_date date,
    funding_end_date date,
    notes text COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT project_grant_funding_pkey PRIMARY KEY (id),
    CONSTRAINT unique_grant_project_section UNIQUE (grant_id, project_id, section_id)
);

ALTER TABLE IF EXISTS public.project_grant_funding
    ADD CONSTRAINT fk_pgf_grant FOREIGN KEY (grant_id)
    REFERENCES public.grants (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.project_grant_funding
    ADD CONSTRAINT fk_pgf_project FOREIGN KEY (project_id)
    REFERENCES public.projects (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.grants
    ADD CONSTRAINT fk_grant_creator FOREIGN KEY (created_by)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;

ALTER TABLE IF EXISTS public.grants
    ADD CONSTRAINT fk_grant_pi FOREIGN KEY (principal_investigator_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_grants_status ON public.grants(status);
CREATE INDEX IF NOT EXISTS idx_grants_type ON public.grants(grant_type);
CREATE INDEX IF NOT EXISTS idx_grants_org ON public.grants(funding_organization);
CREATE INDEX IF NOT EXISTS idx_pgf_grant ON public.project_grant_funding(grant_id);
CREATE INDEX IF NOT EXISTS idx_pgf_project ON public.project_grant_funding(project_id);
