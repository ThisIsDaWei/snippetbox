create table snippets (
    id serial not null primary key,
    title varchar(100) not null,
    content text not null,
    created timestamptz default (now() at time zone 'utc'),
    expires timestamptz not null
);

create index idx_snippets_created on snippets(created);

create table users (
    id serial not null primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    hashed_password char(60) not null,
    created timestamptz default (now() at time zone 'utc')
);

ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:00:00'
);
