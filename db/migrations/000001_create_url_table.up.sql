CREATE TABLE url (
    id SERIAL,
    user_id VARCHAR(255) NOT NULL,
    short_url VARCHAR(255) NULL,
    long_url VARCHAR(255) NULL,
    created_at timestamp NOT NULL default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    CONSTRAINT url_pk PRIMARY KEY (id)
);
CREATE INDEX deleted_at_on_short_url ON url USING btree (deleted_at, short_url);
