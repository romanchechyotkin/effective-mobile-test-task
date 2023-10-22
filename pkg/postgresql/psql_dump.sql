--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3
-- Dumped by pg_dump version 15.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: gender; Type: TYPE; Schema: public; Owner: postgres

--

CREATE TYPE public.gender AS ENUM (
    'male',
    'female'
);


ALTER TYPE public.gender OWNER TO postgres
;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres

--

CREATE TABLE public.users (
    id integer NOT NULL,
    last_name text NOT NULL,
    first_name text NOT NULL,
    second_name text,
    age integer NOT NULL,
    gender public.gender NOT NULL,
    nationality text NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO postgres
;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres

--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres
;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres

--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres

--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres

--

COPY public.users (id, last_name, first_name, second_name, age, gender, nationality, created_at) FROM stdin;
1	Chechyotkin	Roman	Alexandrovich	55	male	CZ	2023-10-22 14:46:54.958413
3	Neronskiy	Lesha	Vyacheslavovich	45	female	BY	2023-10-22 15:03:10.375839
6	Kozlova	Anna	Alexandrovna	51	female	PL	2023-10-22 17:56:55.775437
2	Kostuk	Olga		55	female	UA	2023-10-22 14:47:36.907376
5	asdasda	Nastya	germanovna	5	female	UA	2023-10-22 15:05:56.08489
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres

--

SELECT pg_catalog.setval('public.users_id_seq', 6, true);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres

--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

