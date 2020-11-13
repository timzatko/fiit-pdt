ALTER TABLE tweets
    ADD CONSTRAINT fk_tweets_tweets
        FOREIGN KEY (parent_id) REFERENCES tweets(id);

ALTER TABLE tweets
    ADD CONSTRAINT fk_accounts_tweets
        FOREIGN KEY (author_id) REFERENCES accounts(id);

ALTER TABLE tweets
    ADD CONSTRAINT fk_countries_tweets
        FOREIGN KEY (country_id) REFERENCES countries(id);

ALTER TABLE tweet_hashtags
    ADD CONSTRAINT fk_tweet_hashtags_tweets
        FOREIGN KEY (tweet_id) REFERENCES tweets(id);

ALTER TABLE tweet_hashtags
    ADD CONSTRAINT fk_tweet_hashtags_hashtags
        FOREIGN KEY (hashtag_id) REFERENCES hashtags(id);

ALTER TABLE tweet_mentions
    ADD CONSTRAINT fk_tweet_mentions_tweets
        FOREIGN KEY (tweet_id) REFERENCES tweets(id);

ALTER TABLE tweet_mentions
    ADD CONSTRAINT fk_tweet_mentions_accounts
        FOREIGN KEY (account_id) REFERENCES accounts(id);