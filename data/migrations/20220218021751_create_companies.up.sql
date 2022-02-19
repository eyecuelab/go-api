create extension if not exists pgcrypto;
create extension if not exists citext;

create table companies (
  id serial primary key,
  name text not null,
  slug citext not null,
  description text,
  created_at timestamp without time zone not null default (current_timestamp at time zone 'utc'),
  updated_at timestamp without time zone not null default (current_timestamp at time zone 'utc'),
  deleted_at timestamp without time zone default (current_timestamp at time zone 'utc')
);
create unique index if not exists idx_companies_name on companies ((lower(name)));
create unique index if not exists idx_companies_slug on companies ((lower(slug)));
