create table "user" (
    user_id uuid primary key default uuid_generate_v4(),
    username citext unique not null,
    email citext unique not null,
    password_hash text not null,
    bio text not null default '',
    image text,
    created_at timestamptz not null default now(),
    updated_at timestamptz
);

create trigger set_updated_at
    before update on "user"
    for each row
    execute procedure set_updated_at();

create table "user_token" (
    user_token_id bigserial primary key,
    user_id uuid not null references "user" (user_id) on delete cascade,
    token text not null,
    context text not null,
    created_at timestamptz not null default now()
);

create index on "user_token" (user_id);
create unique index on "user_token" (token, context);
