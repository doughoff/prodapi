begin;
drop type if exists unit;
create type unit as enum ('KG', 'L', 'UNITS', 'OTHER');

create table "products"
(
    "id"                uuid primary key not null default uuid_generate_v4(),
    "status"            status           not null default 'ACTIVE',
    "name"              text             not null,
    "barcode"           text             not null unique,
    "unit"              unit             not null default 'UNITS',
    "batch_control"     boolean          not null default false,
    "conversion_factor" bigint           not null default 1000,
    "created_at"        timestamp        not null default now(),
    "updated_at"        timestamp        not null default now()
);

create index "products_status" on "products" ("status");
commit;