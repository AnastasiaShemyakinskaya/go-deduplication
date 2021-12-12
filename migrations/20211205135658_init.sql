-- +goose Up
-- +goose StatementBegin
create type HashFunction as enum ('MD5', 'SHA128', 'SHA256', 'SHA512');

create table if not exists file_hash
(
    id serial not null primary key,
    hash_id bigint not null,
    file_id bigint not null,
    repeat bigint not null,
    position int[] not null default '{}'
);

create unique index uindex_hash_id_file_id on file_hash (hash_id, file_id);

create table if not exists file
(
    id serial not null primary key,
    hash_function HashFunction not null,
    byte_size integer not null,
    file_name text not null
);

create unique index uindex_file on file (file_name);

create table if not exists hash
(
    id serial not null primary key,
    hash_string bytea not null
);

create unique index uindex_hash_string on hash (hash_string);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists uindex_file;
drop table if exists hash;
drop index if exists uindex_hash_id_file_id;
drop table if exists file;
drop index if exists uindex_hash_string;
drop table if exists file_hash;
drop type if exists HashFunction;
-- +goose StatementEnd
