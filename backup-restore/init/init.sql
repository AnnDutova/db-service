create table if not exists "project" (
    id serial primary key,
    title text not null unique
);