-- liquibase formatted sql
-- changeset author:vlad

CREATE TYPE container_status AS ENUM('online', 'offline');

-- rollback DROP TYPE container_status;