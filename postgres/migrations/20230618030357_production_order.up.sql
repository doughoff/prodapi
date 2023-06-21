begin;

drop type if exists production_step;
create type production_step as enum (
    'PENDING',
    'IN_PROGRESS',
    'COMPLETED'
    );

create table if not exists production_orders
(
    id                   uuid primary key not null default uuid_generate_v4(),
    status               status           not null default 'ACTIVE',
    production_step      production_step  not null default 'PENDING',
    code                 text             not null,
    cycles               bigint           not null,
    recipe_id            uuid             not null,
    created_by_user_id   UUID             NOT NULL,
    cancelled_by_user_id UUID,
    created_at           TIMESTAMP        NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMP        NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_production_order_created_by_user
        FOREIGN KEY (created_by_user_id)
            REFERENCES users (id),
    CONSTRAINT fk_production_order_cancelled_by_user
        FOREIGN KEY (cancelled_by_user_id)
            REFERENCES users (id),
    CONSTRAINT fk_production_order_recipe_id
        FOREIGN KEY (recipe_id)
            REFERENCES recipes (recipe_id)
);

CREATE INDEX production_order_status ON production_orders (status);

create table production_order_cycles
(
    id                  uuid primary key not null default uuid_generate_v4(),
    factor              bigint           not null default 1,
    production_order_id uuid             not null,
    production_step     production_step  not null default 'PENDING',
    cycle_order         bigint           not null,
    completed_at        timestamp,
    constraint fk_cycle_production_order
        foreign key (production_order_id)
            references production_orders (id)
);

create table order_cycles_movements
(
    id          uuid primary key not null default uuid_generate_v4(),
    cycle_id    uuid             not null,
    movement_id uuid             not null,
    constraint fk_cycle_movement_cycle
        foreign key (cycle_id)
            references production_order_cycles (id),
    constraint fk_cycle_movement_movement
        foreign key (movement_id)
            references stock_movements (id)
);

commit;