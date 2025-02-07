FROM liquibase/liquibase:latest

COPY ./migrations /liquibase/migrations
COPY liquibase.properties /liquibase/

CMD ["--defaults-file=/liquibase/liquibase.properties", "update"]