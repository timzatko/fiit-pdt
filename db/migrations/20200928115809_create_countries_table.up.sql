CREATE TABLE countries (
   id				SERIAL PRIMARY KEY,
   code				varchar(2),
   name 			varchar(200),
   UNIQUE (code)
);
