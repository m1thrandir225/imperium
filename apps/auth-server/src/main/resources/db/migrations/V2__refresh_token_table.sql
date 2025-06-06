CREATE TABLE public.refresh_tokens (
    id uuid primary key not null,
    token varchar(255) not null,
    expires_at timestamptz not null,
    user_id uuid constraint refresh_tokens_users_fkey references public.users on delete cascade on update cascade not null
);