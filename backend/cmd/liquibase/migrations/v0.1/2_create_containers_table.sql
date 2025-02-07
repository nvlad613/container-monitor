-- liquibase formatted sql
-- changeset author:vlad

CREATE TABLE containers (
    id SERIAL               PRIMARY KEY,
    docker_id               CHAR(64) NOT NULL UNIQUE,
    name                    VARCHAR(128) NOT NULL,
    status container_status DEFAULT 'offline',
    last_check              TIMESTAMP,
    last_active             TIMESTAMP
);

-- rollback DROP TABLE containers;