CREATE TABLE IF NOT EXISTS public.users (
    id uuid primary key not null,
    email varchar(255) unique not null,
    name varchar(255) not null,
    password varchar(255) not null,
    created_at timestamp(6) not null,
    updated_at timestamp(6)
);

CREATE TABLE IF NOT EXISTS public.clients(
    id uuid primary key not null,
    owner_id uuid constraint clients_users_fkey references public.users on delete cascade on update cascade not null,
    ip_address varchar(255) not null unique,
    name varchar(255) not null unique
);


CREATE TABLE IF NOT EXISTS public.hosts (
    id uuid primary key not null,
    status smallint not null unique constraint hosts_status_check check ( (status >= 0) AND (status <= 2)),
    owner_id uuid constraint hosts_users_fkey references public.users on delete cascade  on update  cascade  not null,
    ip_address varchar(255) not null unique,
    port integer not null unique
);

CREATE TABLE IF NOT EXISTS public.users_hosts(
    hosts_id uuid not null unique constraint users_hosts_hosts_fkey references public.hosts,
    user_id uuid not null constraint users_hosts_users_fkey references public.users
);

CREATE TABLE IF NOT EXISTS public.users_clients (
    clients_id uuid not null unique constraint users_clients_clients_fkey references public.clients,
    user_id uuid not null constraint users_clients_users_fkey references public.users
)
