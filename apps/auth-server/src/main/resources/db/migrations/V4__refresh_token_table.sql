CREATE TABLE public.refresh_tokens
(
    id         uuid primary key                                                                          not null,
    token      varchar(255)                                                                              not null,
    expires_at timestamptz                                                                               not null,
    user_id    uuid
        constraint refresh_tokens_users_fkey references public.users on delete cascade on update cascade not null
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);