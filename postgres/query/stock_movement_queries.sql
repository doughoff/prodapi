-- name: GetStockMovements :many
SELECT count(*) over () as full_count,
       si.id,
       si.status,
       si.type,
       si.date,
       si.entity_id,
       si.document_number,
       e.name           as entity_name,
       si.created_by_user_id,
       u.name           as create_by_user_name,
       si.cancelled_by_user_id,
       uc.name          as cancelled_by_user_name
from stock_movements si
         left join entities e on si.entity_id = e.id
         left join users uc on si.cancelled_by_user_id = uc.id
         left join users u on si.created_by_user_id = u.id
where si.status = any (@status_options::status[])
  and (
        e.name ilike '%' || @search || '%'
      or si.document_number ilike '%' || @search || '%'
    )
  and si.date >= @start_date
order by si.date desc
limit @page_limit offset @page_offset;

-- name: GetStockMovementsByIDS :many
SELECT si.id,
       si.status,
       si.type,
       si.date,
       si.entity_id,
       si.document_number,
       e.name  as entity_name,
       si.created_by_user_id,
       u.name  as create_by_user_name,
       si.cancelled_by_user_id,
       uc.name as cancelled_by_user_name
from stock_movements si
         left join entities e on si.entity_id = e.id
         left join users uc on si.cancelled_by_user_id = uc.id
         left join users u on si.created_by_user_id = u.id
where si.id = any(@stock_movement_id::uuid[]);

-- name: GetStockMovementByID :one
SELECT si.id,
       si.status,
       si.type,
       si.date,
       si.entity_id,
       si.document_number,
       e.name  as entity_name,
       si.created_by_user_id,
       u.name  as create_by_user_name,
       si.cancelled_by_user_id,
       uc.name as cancelled_by_user_name
from stock_movements si
         left join entities e on si.entity_id = e.id
         left join users uc on si.cancelled_by_user_id = uc.id
         left join users u on si.created_by_user_id = u.id
where si.id = @stock_movement_id;

-- name: GetStockMovementItems :many
select smi.id,
       smi.stock_movement_id,
       smi.product_id,
       p.name as product_name,
       smi.quantity,
       smi.price,
       smi.batch,
       smi.created_at,
       smi.updated_at
from stock_movement_items smi
         left join products p on smi.product_id = p.id
where smi.stock_movement_id = any (@stock_movement_ids::uuid[]);

-- name: CreateStockMovement :one
insert into stock_movements(
                            type,
                            date,
                            entity_id,
                            created_by_user_id,
                            document_number
                            )
values (@type, @date, @entity_id, @created_by_user_id, @document_number)
returning id;

-- name: CreateStockMovementItems :copyfrom
insert into stock_movement_items(stock_movement_id, product_id, quantity, price, batch)
values ($1, $2, $3, $4, $5);

-- name: UpdateStockMovement :exec
update stock_movements
 set status = @status,
     entity_id = @entity_id,
     date = @date,
     document_number = @document_number
where id = @id;
