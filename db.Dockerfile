FROM postgres:16.1

WORKDIR /app

# to run scripts after startup, copy *.sql, *.sql.gz, 
# or *.sh scripts to /docker-entrypoint-initdb.d/ 
COPY init_db1.sql /docker-entrypoint-initdb.d/
COPY init_db2.sh /docker-entrypoint-initdb.d/

