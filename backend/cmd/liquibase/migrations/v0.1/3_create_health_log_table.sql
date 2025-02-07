-- liquibase formatted sql
-- changeset author:vlad

CREATE TABLE health_log (
    container_id            INTEGER NOT NULL REFERENCES containers (id),
    status container_status NOT NULL DEFAULT 'offline',
    timestamp               TIMESTAMP NOT NULL,
    PRIMARY KEY(container_id, timestamp)
);

-- rollback DROP TABLE containers;