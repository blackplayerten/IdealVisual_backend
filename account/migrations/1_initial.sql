-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS account (
    id serial PRIMARY KEY,
    email varchar(64) UNIQUE NOT NULL,
    password varchar(64) NOT NULL,

    username varchar(64) UNIQUE NOT NULL,
    avatar varchar(100)
);

CREATE TABLE IF NOT EXISTS post (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    acc integer REFERENCES account NOT NULL,
    photo varchar(100) NOT NULL,
    photo_index integer DEFAULT 0,
    date timestamp with time zone,
    place varchar(256),
    text varchar(5000)
);

-- +migrate Down

DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS account;

DROP EXTENSION IF EXISTS "uuid-ossp";
