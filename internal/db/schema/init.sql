-- Simulate IF NOT EXIST on CREATE DATABASE
SELECT 'CREATE DATABASE timeline'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'timeline')\gexec

create table if not exists types
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name       text,
    color      text
);

alter table types
    owner to postgres;

create index if not exists idx_types_deleted_at
    on types (deleted_at);

INSERT INTO types(id, created_at, updated_at, deleted_at, name, color)
VALUES (1, now(), now(), null, 'normal', 'white')
