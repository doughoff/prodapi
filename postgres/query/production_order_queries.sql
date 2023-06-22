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

-- name: GetOrderCyclesMovements :many
select ocm.id,
       ocm.cycle_id,
       ocm.movement_id,
       coalesce(ri.recipe_id, smi.product_id) as product_id,
       coalesce(smi.quantity, 0) as cantidad,
       coalesce(smi.price, 0) as price
from order_cycles_movements ocm
         join production_order_cycles poc on ocm.cycle_id = poc.id
         left join production_orders po on poc.production_order_id = po.id
         left join recipe_ingredients ri on po.recipe_id = ri.recipe_id
         left join stock_movement_items smi on ocm.movement_id = smi.stock_movement_id
where ocm.cycle_id = any (@cycle_ids::uuid[]);

-- name: GetProductionOrderMovements :many
select *
from production_order_cycles poc
    join order_cycles_movements ocm on poc.id = ocm.cycle_id
    join stock_movements sm on ocm.movement_id = sm.id
where poc.production_order_id = @production_order_id;
