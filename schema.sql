--
-- PostgreSQL database dump
--

\restrict 9ulNUgjPOftAELG2QhGr54o7h8bC3hG7siG3e9WmFGQ6E4S2ZZDJGJKoCdP4W4n

-- Dumped from database version 17.7 (Ubuntu 17.7-0ubuntu0.25.04.1)
-- Dumped by pg_dump version 17.7 (Ubuntu 17.7-0ubuntu0.25.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: contacts; Type: TABLE; Schema: common; Owner: conor
--

CREATE TABLE common.contacts (
    id integer NOT NULL,
    name text NOT NULL,
    phone text NOT NULL,
    email text NOT NULL,
    customer_id integer,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE common.contacts OWNER TO conor;

--
-- Name: contacts_id_seq; Type: SEQUENCE; Schema: common; Owner: conor
--

CREATE SEQUENCE common.contacts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE common.contacts_id_seq OWNER TO conor;

--
-- Name: contacts_id_seq; Type: SEQUENCE OWNED BY; Schema: common; Owner: conor
--

ALTER SEQUENCE common.contacts_id_seq OWNED BY common.contacts.id;


--
-- Name: notes; Type: TABLE; Schema: common; Owner: conor
--

CREATE TABLE common.notes (
    id integer NOT NULL,
    content text NOT NULL,
    customer_id integer,
    job_id integer,
    stairway_id integer,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE common.notes OWNER TO conor;

--
-- Name: notes_id_seq; Type: SEQUENCE; Schema: common; Owner: conor
--

CREATE SEQUENCE common.notes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE common.notes_id_seq OWNER TO conor;

--
-- Name: notes_id_seq; Type: SEQUENCE OWNED BY; Schema: common; Owner: conor
--

ALTER SEQUENCE common.notes_id_seq OWNED BY common.notes.id;


--
-- Name: brackets; Type: TABLE; Schema: handrails; Owner: conor
--

CREATE TABLE handrails.brackets (
    id integer NOT NULL,
    name character varying(255),
    nickname character varying(255),
    supplier_id integer,
    price double precision,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE handrails.brackets OWNER TO conor;

--
-- Name: brackets_id_seq; Type: SEQUENCE; Schema: handrails; Owner: conor
--

CREATE SEQUENCE handrails.brackets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE handrails.brackets_id_seq OWNER TO conor;

--
-- Name: brackets_id_seq; Type: SEQUENCE OWNED BY; Schema: handrails; Owner: conor
--

ALTER SEQUENCE handrails.brackets_id_seq OWNED BY handrails.brackets.id;


--
-- Name: materials; Type: TABLE; Schema: materials; Owner: conor
--

CREATE TABLE materials.materials (
    id integer NOT NULL,
    name text NOT NULL,
    nickname text DEFAULT ''::text NOT NULL,
    treatment text DEFAULT 'UT'::text NOT NULL,
    blurb text DEFAULT ''::text NOT NULL,
    units text DEFAULT 'mm'::text NOT NULL,
    length double precision DEFAULT 2400 NOT NULL,
    width double precision DEFAULT 1200 NOT NULL,
    thickness double precision DEFAULT 20 NOT NULL,
    max_span double precision DEFAULT 500 NOT NULL,
    density double precision DEFAULT 500 NOT NULL,
    max_overhang double precision DEFAULT 10 NOT NULL,
    radius double precision DEFAULT 3 NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE materials.materials OWNER TO conor;

--
-- Name: materials_id_seq; Type: SEQUENCE; Schema: materials; Owner: conor
--

CREATE SEQUENCE materials.materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE materials.materials_id_seq OWNER TO conor;

--
-- Name: materials_id_seq; Type: SEQUENCE OWNED BY; Schema: materials; Owner: conor
--

ALTER SEQUENCE materials.materials_id_seq OWNED BY materials.materials.id;


--
-- Name: materials_notes; Type: TABLE; Schema: materials; Owner: conor
--

CREATE TABLE materials.materials_notes (
    id integer NOT NULL,
    material_id integer,
    note_id integer,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE materials.materials_notes OWNER TO conor;

--
-- Name: materials_notes_id_seq; Type: SEQUENCE; Schema: materials; Owner: conor
--

CREATE SEQUENCE materials.materials_notes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE materials.materials_notes_id_seq OWNER TO conor;

--
-- Name: materials_notes_id_seq; Type: SEQUENCE OWNED BY; Schema: materials; Owner: conor
--

ALTER SEQUENCE materials.materials_notes_id_seq OWNED BY materials.materials_notes.id;


--
-- Name: suppliers; Type: TABLE; Schema: materials; Owner: conor
--

CREATE TABLE materials.suppliers (
    id integer NOT NULL,
    name text NOT NULL,
    website text NOT NULL,
    contact_id integer NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE materials.suppliers OWNER TO conor;

--
-- Name: suppliers_id_seq; Type: SEQUENCE; Schema: materials; Owner: conor
--

CREATE SEQUENCE materials.suppliers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE materials.suppliers_id_seq OWNER TO conor;

--
-- Name: suppliers_id_seq; Type: SEQUENCE OWNED BY; Schema: materials; Owner: conor
--

ALTER SEQUENCE materials.suppliers_id_seq OWNED BY materials.suppliers.id;


--
-- Name: suppliers_to_materials; Type: TABLE; Schema: materials; Owner: conor
--

CREATE TABLE materials.suppliers_to_materials (
    id integer NOT NULL,
    priority integer NOT NULL,
    supplier_id integer NOT NULL,
    material_id integer NOT NULL,
    price double precision NOT NULL,
    lead_time interval DEFAULT '14 days'::interval NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE materials.suppliers_to_materials OWNER TO conor;

--
-- Name: suppliers_to_materials_id_seq; Type: SEQUENCE; Schema: materials; Owner: conor
--

CREATE SEQUENCE materials.suppliers_to_materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE materials.suppliers_to_materials_id_seq OWNER TO conor;

--
-- Name: suppliers_to_materials_id_seq; Type: SEQUENCE OWNED BY; Schema: materials; Owner: conor
--

ALTER SEQUENCE materials.suppliers_to_materials_id_seq OWNED BY materials.suppliers_to_materials.id;


--
-- Name: uses; Type: TABLE; Schema: materials; Owner: conor
--

CREATE TABLE materials.uses (
    id integer NOT NULL,
    material_id integer NOT NULL,
    exterior boolean DEFAULT false NOT NULL,
    interior boolean DEFAULT true NOT NULL,
    early_access boolean DEFAULT false NOT NULL,
    stringers boolean DEFAULT false NOT NULL,
    risers boolean DEFAULT true NOT NULL,
    treads boolean DEFAULT true NOT NULL,
    timber boolean DEFAULT false NOT NULL,
    panel boolean DEFAULT false NOT NULL,
    published boolean DEFAULT false NOT NULL,
    archived boolean DEFAULT false NOT NULL,
    handrail boolean DEFAULT false NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE materials.uses OWNER TO conor;

--
-- Name: uses_id_seq; Type: SEQUENCE; Schema: materials; Owner: conor
--

CREATE SEQUENCE materials.uses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE materials.uses_id_seq OWNER TO conor;

--
-- Name: uses_id_seq; Type: SEQUENCE OWNED BY; Schema: materials; Owner: conor
--

ALTER SEQUENCE materials.uses_id_seq OWNED BY materials.uses.id;


--
-- Name: contacts id; Type: DEFAULT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.contacts ALTER COLUMN id SET DEFAULT nextval('common.contacts_id_seq'::regclass);


--
-- Name: notes id; Type: DEFAULT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.notes ALTER COLUMN id SET DEFAULT nextval('common.notes_id_seq'::regclass);


--
-- Name: brackets id; Type: DEFAULT; Schema: handrails; Owner: conor
--

ALTER TABLE ONLY handrails.brackets ALTER COLUMN id SET DEFAULT nextval('handrails.brackets_id_seq'::regclass);


--
-- Name: materials id; Type: DEFAULT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.materials ALTER COLUMN id SET DEFAULT nextval('materials.materials_id_seq'::regclass);


--
-- Name: materials_notes id; Type: DEFAULT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.materials_notes ALTER COLUMN id SET DEFAULT nextval('materials.materials_notes_id_seq'::regclass);


--
-- Name: suppliers id; Type: DEFAULT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.suppliers ALTER COLUMN id SET DEFAULT nextval('materials.suppliers_id_seq'::regclass);


--
-- Name: suppliers_to_materials id; Type: DEFAULT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.suppliers_to_materials ALTER COLUMN id SET DEFAULT nextval('materials.suppliers_to_materials_id_seq'::regclass);


--
-- Name: uses id; Type: DEFAULT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.uses ALTER COLUMN id SET DEFAULT nextval('materials.uses_id_seq'::regclass);


--
-- Name: contacts contacts_pkey; Type: CONSTRAINT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.contacts
    ADD CONSTRAINT contacts_pkey PRIMARY KEY (id);


--
-- Name: notes notes_pkey; Type: CONSTRAINT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.notes
    ADD CONSTRAINT notes_pkey PRIMARY KEY (id);


--
-- Name: brackets brackets_pkey; Type: CONSTRAINT; Schema: handrails; Owner: conor
--

ALTER TABLE ONLY handrails.brackets
    ADD CONSTRAINT brackets_pkey PRIMARY KEY (id);


--
-- Name: materials_notes materials_notes_pkey; Type: CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.materials_notes
    ADD CONSTRAINT materials_notes_pkey PRIMARY KEY (id);


--
-- Name: materials materials_pkey; Type: CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.materials
    ADD CONSTRAINT materials_pkey PRIMARY KEY (id);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- Name: suppliers_to_materials suppliers_to_materials_pkey; Type: CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.suppliers_to_materials
    ADD CONSTRAINT suppliers_to_materials_pkey PRIMARY KEY (id);


--
-- Name: uses uses_pkey; Type: CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.uses
    ADD CONSTRAINT uses_pkey PRIMARY KEY (id);


--
-- Name: contacts contacts_customer_id_fkey; Type: FK CONSTRAINT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.contacts
    ADD CONSTRAINT contacts_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customers.customers(id);


--
-- Name: notes notes_customer_id_fkey; Type: FK CONSTRAINT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.notes
    ADD CONSTRAINT notes_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customers.customers(id);


--
-- Name: notes notes_job_id_fkey; Type: FK CONSTRAINT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.notes
    ADD CONSTRAINT notes_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs.jobs(id);


--
-- Name: notes notes_stairway_id_fkey; Type: FK CONSTRAINT; Schema: common; Owner: conor
--

ALTER TABLE ONLY common.notes
    ADD CONSTRAINT notes_stairway_id_fkey FOREIGN KEY (stairway_id) REFERENCES stairways.stairways(id);


--
-- Name: brackets brackets_supplier_id_fkey; Type: FK CONSTRAINT; Schema: handrails; Owner: conor
--

ALTER TABLE ONLY handrails.brackets
    ADD CONSTRAINT brackets_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES materials.suppliers(id);


--
-- Name: materials_notes materials_notes_material_id_fkey; Type: FK CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.materials_notes
    ADD CONSTRAINT materials_notes_material_id_fkey FOREIGN KEY (material_id) REFERENCES materials.materials(id);


--
-- Name: materials_notes materials_notes_note_id_fkey; Type: FK CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.materials_notes
    ADD CONSTRAINT materials_notes_note_id_fkey FOREIGN KEY (note_id) REFERENCES common.notes(id);


--
-- Name: suppliers suppliers_contact_id_fkey; Type: FK CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.suppliers
    ADD CONSTRAINT suppliers_contact_id_fkey FOREIGN KEY (contact_id) REFERENCES common.contacts(id);


--
-- Name: suppliers_to_materials suppliers_to_materials_material_id_fkey; Type: FK CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.suppliers_to_materials
    ADD CONSTRAINT suppliers_to_materials_material_id_fkey FOREIGN KEY (material_id) REFERENCES materials.materials(id);


--
-- Name: suppliers_to_materials suppliers_to_materials_supplier_id_fkey; Type: FK CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.suppliers_to_materials
    ADD CONSTRAINT suppliers_to_materials_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES materials.suppliers(id);


--
-- Name: uses uses_material_id_fkey; Type: FK CONSTRAINT; Schema: materials; Owner: conor
--

ALTER TABLE ONLY materials.uses
    ADD CONSTRAINT uses_material_id_fkey FOREIGN KEY (material_id) REFERENCES materials.materials(id);


--
-- PostgreSQL database dump complete
--

\unrestrict 9ulNUgjPOftAELG2QhGr54o7h8bC3hG7siG3e9WmFGQ6E4S2ZZDJGJKoCdP4W4n

