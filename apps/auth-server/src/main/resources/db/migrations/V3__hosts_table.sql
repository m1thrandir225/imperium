CREATE TYPE host_status_enum AS ENUM ('OFFLINE', 'AVAILABLE', 'INUSE', 'DISABLED');

CREATE TABLE IF NOT EXISTS public.hosts
(
    id         uuid primary key                                                                 not null,
    name       text                                                                             not null,
    status host_status_enum not null DEFAULT 'AVAILABLE',
    owner_id   uuid
        constraint hosts_users_fkey references public.users on delete cascade on update cascade not null,
    ip_address varchar(255)                                                                     not null,
    port       integer                                                                          not null,

    CONSTRAINT uk_host_name_ip_port_status UNIQUE (name, ip_address, port, status)
);

CREATE INDEX IF NOT EXISTS idx_hosts_owner_id ON hosts (owner_id);
CREATE INDEX IF NOT EXISTS idx_hosts_status ON hosts (status);
CREATE INDEX IF NOT EXISTS idx_hosts_name_ip_port ON hosts (name, ip_address, port);