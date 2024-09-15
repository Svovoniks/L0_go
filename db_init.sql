create database orders;

\connect orders

create table "order" (
    id text primary key,
    json_data text
);
