create table if not exists users (
  id serial primary key,
  first_name text not null,
  last_name text not null,
  email citext not null,
  roles text[] not null default '{}'::text[],
  reset_token text,
  last_signin_at timestamp without time zone,
  company_id int,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);
create unique index if not exists idx_users_email on users ((lower(email)));

create type credential_source as enum ('password', 'confirm_email', 'password_reset');
create table if not exists credentials (
  id serial primary key,
  user_id int not null references users(id),
  source credential_source not null default('password'),
  value text not null,
  expires_at timestamp without time zone,
  created_at timestamp without time zone not null default (current_timestamp at time zone 'utc'),
  updated_at timestamp without time zone not null default (current_timestamp at time zone 'utc')
);
create unique index if not exists idx_credentials on credentials (user_id, source);
