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
    ip_address varchar(255) not null,
    name varchar(255) not null,

    CONSTRAINT uk_client_name_ip_address UNIQUE (name, ip_address)


);


CREATE TABLE IF NOT EXISTS public.hosts (
    id uuid primary key not null,
    name text not null,
    status smallint not null constraint hosts_status_check check ( (status >= 0) AND (status <= 2)),
    owner_id uuid constraint hosts_users_fkey references public.users on delete cascade  on update  cascade  not null,
    ip_address varchar(255) not null,
    port integer not null,

    CONSTRAINT uk_host_name_ip_port_status UNIQUE (name, ip_address, port, status)
);