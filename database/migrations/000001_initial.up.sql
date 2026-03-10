begin;
create extension if not exists "pgcrypto";
create type asset_type as enum (
    'laptop',
    'keyboard',
    'mouse',
    'mobile'
    );

create type asset_status as enum (
    'available',
    'assigned',
    'in_service',
    'for_repair',
    'damaged'
    );

create type user_role as enum (
    'admin',
    'employee',
    'project-manager',
    'asset-manager',
    'employee-manager'
    );

create type user_type as enum(
     'full-time',
     'intern',
     'freelancer'
     );


create type owner_type as enum (
    'client',
    'company'
    );

create table if not exists users
(
    id            uuid primary key default gen_random_uuid(),
    name          text not null,
    email         text not null,
    role          user_role        default 'employee',
    type          user_type not null,
    phone_no      text not null,
    password_hash text not null,
    joining_date  date not null ,
    created_at    timestamp        default current_timestamp,
    archived_at   timestamptz
    );

create unique index idx_unique_email on users (email) where archived_at is NULL;



create table assets
(
    id             uuid primary key default gen_random_uuid(),
    brand          text             not null,
    model          text             not null,
    serial_no      text unique      not null,
    type           asset_type       not null,
    status         asset_status     default 'available',
    owner          owner_type       not null,

    assigned_by_id uuid references users(id),
    assigned_to    uuid references users (id),
    assigned_on    timestamptz,

    warranty_start timestamptz      not null,
    warranty_end   timestamptz      not null,

    service_start  timestamptz,
    service_end    timestamptz,
    returned_on    timestamptz,

    created_at     timestamptz      default now(),
    updated_at timestamptz,
    archived_at timestamptz,
    archived_by uuid references users(id)
);

create table if not exists user_session
(
    id          uuid primary key default gen_random_uuid(),
    user_id     uuid references users (id) NOT NULL,
    created_at  timestamp        default current_timestamp,
    archived_at timestamptz
    );

create table laptop
(
    id        uuid primary key default gen_random_uuid(),
    asset_id  uuid unique references assets (id),
    processor text,
    ram       text,
    storage   text,
    os        text,
    charger   text,
    password  text             not null

);

create type connection_type as enum ('wired', 'wireless');

create table keyboard
(
    id uuid primary key default gen_random_uuid(),
    asset_id     uuid unique references assets (id),
    layout       text,
    connectivity connection_type
);

create table mouse
(
    id           uuid primary key default gen_random_uuid(),
    asset_id     uuid unique references assets (id),
    dpi          int,
    connectivity connection_type
);

create table mobile
(
    id       uuid primary key default gen_random_uuid(),
    asset_id uuid unique references assets(id),
    os       text             not null,
    ram      text             not null,
    storage  text             not null,
    charger  text,
    password text             not null
);

commit;