CREATE TABLE hashtags (
    id				SERIAL PRIMARY KEY,
    value			text,
    UNIQUE (value)
)