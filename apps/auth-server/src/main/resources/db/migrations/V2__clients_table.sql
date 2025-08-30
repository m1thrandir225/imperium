CREATE TABLE IF NOT EXISTS public.clients
(
    id         uuid primary key                                                                   not null,
    owner_id   uuid
        constraint clients_users_fkey references public.users on delete cascade on update cascade not null,
    ip_address varchar(255)                                                                       not null,
    name       varchar(255)                                                                       not null,

    CONSTRAINT uk_client_name_ip_address UNIQUE (name, ip_address)
);

CREATE INDEX IF NOT EXISTS idx_clients_owner_id ON clients(owner_id);
CREATE INDEX IF NOT EXISTS idx_clients_name_ip ON clients(name, ip_address);