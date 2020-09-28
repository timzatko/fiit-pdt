CREATE TABLE tweets (
    id				varchar(20) NOT NULL,
    content			text,
    location 		geometry(point, 4326),
    retweet_count	integer NOT NULL,
    favorite_count	integer NOT NULL,
    happened_at		timestamptz NOT NULL,
    author_id		bigint NOT NULL,
    country_id		integer,
    parent_id		varchar(20),
    PRIMARY KEY (id),
    FOREIGN KEY (parent_id) REFERENCES tweets(id)
)