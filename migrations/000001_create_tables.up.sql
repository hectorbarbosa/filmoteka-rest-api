CREATE TYPE g AS ENUM ('M', 'F');

CREATE TABLE IF NOT EXISTS public.films (
	id SERIAL PRIMARY KEY,
	name varchar(150) NOT NULL,
	description varchar(500),
	release_year smallint NOT NULL,
    rating float 
);

CREATE TABLE IF NOT EXISTS public.actors (
	id SERIAL PRIMARY KEY,
	name varchar(100) NOT NULL,
	gender g NOT NULL,
	birth_date date NOT NULL
);

CREATE UNIQUE INDEX ON public.films(name, release_year);