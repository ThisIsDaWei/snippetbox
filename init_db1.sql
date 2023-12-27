-- Reference for init_db.sh

CREATE DATABASE snippetbox WITH ENCODING 'UTF8' LC_COLLATE = 'en_US.UTF-8' 
LC_CTYPE = 'en_US.UTF-8' TEMPLATE template0;

-- Connect to database
\c snippetbox;

-- Snippets
-- serial is auto-incrementing integer type
create table snippets (
    id serial not null primary key,
    title varchar(100) not null,
    content text not null,
    created timestamptz default (now() at time zone 'utc'),
    expires timestamptz not null
);

create index idx_snippets_created on snippets(created);


-- Sessions
create table sessions (
    token char(43) primary key,
    data bytea not null,
    expiry timestamptz(6) not null
);

create index sessions_expiry_idx on sessions (expiry);


-- Users
create table users (
    id serial not null primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    hashed_password char(60) not null,
    created timestamptz default (now() at time zone 'utc')
);

alter table users add constraint users_email_key unique (email);


-- Dummy records 
INSERT INTO snippets (title, content, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    now() + interval '365 days'
);

INSERT INTO snippets (title, content, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    now() + interval '365 days'
);

INSERT INTO snippets (title, content, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    now() + interval '7 days'
);



-- For unit tests and end-to-end tests.
CREATE DATABASE test_snippetbox WITH ENCODING 'UTF8' LC_COLLATE = 'en_US.UTF-8' 
LC_CTYPE = 'en_US.UTF-8' TEMPLATE template0;

\c test_snippetbox;

create user test_web with login password 'pass';

grant usage on schema public to test_web;

alter default privileges in schema public
    grant all on tables to test_web;

alter default privileges in schema public
    grant usage, select on sequences to test_web;
