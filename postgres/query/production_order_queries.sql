-- name: GetProductionOrders :many
SELECT count(*) over () as full_count,
       po.id,
       po.status,
       po.production_step,
       po.code,
       po.cycles,
       po.recipe_id,
       r.name,
       r.produced_quantity,
       p.name,
       p.unit,
       po.created_by_user_id,
       u.name           as create_by_user_name,
       po.cancelled_by_user_id,
       ud.name          as cancelled_by_user_name,
       po.created_at
from production_orders po
         left join recipes r on po.recipe_id = r.recipe_id
         left join products p on r.product_id = p.id
         left join users u on po.created_by_user_id = u.id
         left join users ud on po.cancelled_by_user_id = ud.id
where po.status = any (@status_optiuons::status[])
  and (
            r.name ilike '%' || @search || '%'
        or p.name ilike '%' || @search || '%'
    )
  and po.created_at >= @start_date
order by po.created_at desc
limit @page_limit offset @page_offset;

-- name: GetProductionOrderByID :one
SELECT po.id,
       po.status,
       po.production_step,
       po.code,
       po.cycles,
       po.recipe_id,
       r.name,
       r.produced_quantity,
       p.name,
       p.unit,
       po.created_by_user_id,
       u.name  as create_by_user_name,
       po.cancelled_by_user_id,
       ud.name as cancelled_by_user_name,
       po.created_at
from production_orders po
         left join recipes r on po.recipe_id = r.recipe_id
         left join products p on r.product_id = p.id
         left join users u on po.created_by_user_id = u.id
         left join users ud on po.cancelled_by_user_id = ud.id
where po.id = @production_order_id;

-- name: GetProductionOrderCycles :many
select poc.id,
       factor,
       production_order_id,
       production_step,
       cycle_order,
       completed_at
from production_order_cycles poc
where poc.production_order_id = any (@production_order_ids::uuid[]);

-- name: GetProductionOrderMovements :many
select ocm.cycle_id    as production_order_cycle_id,
       ocm.movement_id as stock_movement_id
from production_order_cycles poc
         join order_cycles_movements ocm on poc.id = ocm.cycle_id
where poc.production_order_id = @production_order_id;

-- name: CreateProductionOrder :one
insert into production_orders(code, cycles, recipe_id, created_by_user_id)
values (@code, @cycles, @recipe_id, @created_by_user_id)
returning id;

-- name: UpdateProductionOrder :exec
update production_orders
    set status = @status,
        production_step = @production_step
where id = @id;

-- name: CreateProductionOrderCycles :copyfrom
insert into production_order_cycles(factor, production_order_id, cycle_order)
values ($1, $2, $3);

-- name: CreateOrderCycleMovement :one
insert into order_cycles_movements(cycle_id, movement_id)
VALUES (@cycle_id, @movement_id) returning id;