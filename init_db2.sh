#!/bin/bash
set -e # Exit immediately if a command exits with a non-zero status.

# psql variable 
psql -v ON_ERROR_STOP=1 <<-EOSQL

\c snippetbox;

-- Create web user
create user web with login password '${USER_WEB_PASSWORD}';

-- "public" is the default schema
grant usage on schema public to web;

-- Grant permissions after creation of tables, 
-- else permissions are set to default (no permissions). 
grant select, insert, update, delete on all tables in schema public to web;

-- Grant permissions to user to use auto-incrementing columns
grant usage, select on all sequences in schema public to web;

EOSQL