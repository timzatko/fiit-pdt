CREATE TABLE accounts (
    id					bigint NOT NULL PRIMARY KEY,
    screen_name			varchar(200),
    name				varchar(200),
    description		 	text,
    followers_count		integer,
    friends_count		integer,
    statuses_count		integer
);

ALTER TABLE tweets
    ADD CONSTRAINT fk_accounts_tweets
        FOREIGN KEY (author_id) REFERENCES accounts(id);
