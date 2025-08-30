CREATE TABLE IF NOT EXISTS public.users
(
    id         uuid primary key    not null,
    email      varchar(255) unique not null,
    name       varchar(255)        not null,
    password   varchar(255)        not null,
    created_at timestamp(6)        not null,
    updated_at timestamp(6)
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
