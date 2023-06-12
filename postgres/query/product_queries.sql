-- name: GetProducts :many
select count(*) over () as full_count,
       p.*
from products p
where p.status = ANY (@status_options::status[])
  and (
            p.name ilike '%' || @search || '%'
        or p.barcode ilike '%' || @search || '%'
    )
order by created_at
limit @page_limit offset @page_offset;

-- name: GetProductByID :one
select *
from products
where id = @id
limit 1;

-- name: GetProductByBarcode :one
select *
from products
where barcode = @barcode
limit 1;

-- name: CreateProduct :one
insert into products (name, barcode, unit, batch_control, conversion_factor)
values (@name, @barcode, @unit, @batch_control, @conversion_factor)
returning *;

-- name: UpdateProductByID :one
update products
set status            = @status,
    name              = @name,
    barcode           = @barcode,
    unit              = @unit,
    batch_control     = @batch_control,
    conversion_factor = @conversion_factor
where id = @id
returning *;
