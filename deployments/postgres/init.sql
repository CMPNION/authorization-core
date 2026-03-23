create table if not exists users (
  id bigserial primary key,
  login text not null unique,
  password_hash text not null,
  created_at timestamptz not null default now()
);

create index if not exists idx_users_login on users (login);
