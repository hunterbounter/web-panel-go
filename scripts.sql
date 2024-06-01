/*
    Hunter Bounter Database Script v1.0
 */


-- DROP SCHEMA public;

CREATE SCHEMA public AUTHORIZATION pg_database_owner;

COMMENT ON SCHEMA public IS 'standard public schema';

-- DROP TYPE public.openvas_results;

CREATE TYPE public.openvas_results AS (
    id SERIAL,
    report_id varchar(255),
    host text,
    port int4,
    nvt_name text,
    nvt_oid text,
    severity numeric,
    description text,
    solution text,
    created_at timestamp);

-- DROP TYPE public.scan_results;

CREATE TYPE public.scan_results AS (
    id SERIAL,
    alert_id int4,
    url text,
    risk text,
    description text,
    solution text,
    other_info text,
    reference text,
    cwe_id int4,
    wasc_id int4,
    machine_id varchar);

-- DROP TYPE public.targets;

CREATE TYPE public.targets AS (
    id SERIAL,
    value varchar,
    "type" int4,
    status int4);

-- DROP TYPE public."_openvas_results";

CREATE TYPE public."_openvas_results" (
	INPUT = array_in,
	OUTPUT = array_out,
	RECEIVE = array_recv,
	SEND = array_send,
	ANALYZE = array_typanalyze,
	ALIGNMENT = 8,
	STORAGE = any,
	CATEGORY = A,
	ELEMENT = public.openvas_results,
	DELIMITER = ',');

-- DROP TYPE public."_scan_results";

CREATE TYPE public."_scan_results" (
	INPUT = array_in,
	OUTPUT = array_out,
	RECEIVE = array_recv,
	SEND = array_send,
	ANALYZE = array_typanalyze,
	ALIGNMENT = 8,
	STORAGE = any,
	CATEGORY = A,
	ELEMENT = public.scan_results,
	DELIMITER = ',');

-- DROP TYPE public."_targets";

CREATE TYPE public."_targets" (
	INPUT = array_in,
	OUTPUT = array_out,
	RECEIVE = array_recv,
	SEND = array_send,
	ANALYZE = array_typanalyze,
	ALIGNMENT = 8,
	STORAGE = any,
	CATEGORY = A,
	ELEMENT = public.targets,
	DELIMITER = ',');

-- DROP SEQUENCE public.openvas_results_id_seq;

CREATE SEQUENCE public.openvas_results_id_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 2147483647
    START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.openvas_results_id_seq OWNER TO zap_user;
GRANT ALL ON SEQUENCE public.openvas_results_id_seq TO zap_user;

-- DROP SEQUENCE public.scan_results_id_seq;

CREATE SEQUENCE public.scan_results_id_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 2147483647
    START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.scan_results_id_seq OWNER TO zap_user;
GRANT ALL ON SEQUENCE public.scan_results_id_seq TO zap_user;

-- DROP SEQUENCE public.targets_id_seq;

CREATE SEQUENCE public.targets_id_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 2147483647
    START 1
	CACHE 1
	NO CYCLE;

-- Permissions

ALTER SEQUENCE public.targets_id_seq OWNER TO zap_user;
GRANT ALL ON SEQUENCE public.targets_id_seq TO zap_user;
-- public.openvas_results definition

-- Drop table

-- DROP TABLE public.openvas_results;

CREATE TABLE public.openvas_results (
                                        id SERIAL NOT NULL,
                                        report_id varchar(255) NOT NULL,
                                        host text NOT NULL,
                                        port int4 NOT NULL,
                                        nvt_name text NOT NULL,
                                        nvt_oid text NOT NULL,
                                        severity numeric NOT NULL,
                                        description text NULL,
                                        solution text NULL,
                                        created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
                                        CONSTRAINT openvas_results_pkey PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE public.openvas_results OWNER TO zap_user;
GRANT ALL ON TABLE public.openvas_results TO zap_user;


-- public.scan_results definition

-- Drop table

-- DROP TABLE public.scan_results;

CREATE TABLE public.scan_results (
                                     id SERIAL NOT NULL,
                                     alert_id int4 NOT NULL,
                                     url text NOT NULL,
                                     risk text NOT NULL,
                                     description text NOT NULL,
                                     solution text NULL,
                                     other_info text NULL,
                                     reference text NULL,
                                     cwe_id int4 NULL,
                                     wasc_id int4 NULL,
                                     machine_id varchar NULL,
                                     CONSTRAINT scan_results_pkey PRIMARY KEY (id),
                                     CONSTRAINT scan_results_unique UNIQUE (id)
);

-- Permissions

ALTER TABLE public.scan_results OWNER TO zap_user;
GRANT ALL ON TABLE public.scan_results TO zap_user;


-- public.targets definition

-- Drop table

-- DROP TABLE public.targets;

CREATE TABLE public.targets (
                                id SERIAL NOT NULL,
                                value varchar NULL,
                                "type" int4 NULL, -- 1 - Domain 2 - Ipv4 Addr
                                status int4 NULL, -- 1 - Waiting 2 - Started 3 - Finish
                                CONSTRAINT targets_unique UNIQUE (id)
);

-- Column comments

COMMENT ON COLUMN public.targets."type" IS '1 - Domain 2 - Ipv4 Addr';
COMMENT ON COLUMN public.targets.status IS '1 - Waiting 2 - Started 3 - Finish';

-- Permissions

ALTER TABLE public.targets OWNER TO zap_user;
GRANT ALL ON TABLE public.targets TO zap_user;




-- Permissions

GRANT ALL ON SCHEMA public TO pg_database_owner;
GRANT USAGE ON SCHEMA public TO public;