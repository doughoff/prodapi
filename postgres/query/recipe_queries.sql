-- name: GetRecipes :many
select count(*) over () as full_count,
       r.recipe_id,
       r.recipe_group_id,
       r.status,
       r.name,
       r.product_id     as product_id,
       p.name           as product_name,
       p.unit           as product_unit,
       r.produced_quantity,
       r.revision,
       r.is_current,
       r.created_by_user_id,
       u.name           as created_by_user_name,
       r.created_at
from recipes r
         join users u on u.id = r.created_by_user_id
         join products p on r.product_id = p.id
where r.status = any (@status_options::status[])
  and is_current = true
  and (
            r.name ilike '%' || @search || '%'
        or p.name ilike '%' || @search || '%'
    )
order by r.created_at desc
limit @page_limit offset @page_offset;

-- name: GetRecipeIngredients :many
select ri.id,
       ri.recipe_id,
       ri.product_id,
       p.name as product_name,
       p.unit as product_unit,
       ri.quantity
from recipe_ingredients ri
         join products p on ri.product_id = p.id
where recipe_id = any (@recipe_ids::uuid[]);

-- name: GetRecipeByID :one
select r.recipe_id,
       r.recipe_group_id,
       r.name,
       r.product_id         as product_id,
       p.name               as product_name,
       p.unit               as product_unit,
       r.produced_quantity,
       r.status,
       r.revision,
       r.is_current,
       r.created_by_user_id,
       u.name as created_by_user_name,
       r.created_at
from recipes r
         join users u on u.id = r.created_by_user_id
         join products p on r.product_id = p.id
where r.recipe_id = @recipe_id
limit 1;

-- name: GetRecipesByGroupID :many
select r.recipe_id,
       r.recipe_group_id,
       r.name,
       r.product_id as product_id,
       p.name       as product_name,
       p.unit       as product_unit,
       r.produced_quantity,
       r.status,
       r.revision,
       r.is_current,
       r.created_by_user_id,
       u.name       as created_by_user_name,
       r.created_at
from recipes r
         join users u on u.id = r.created_by_user_id
         join products p on r.product_id = p.id
where r.recipe_group_id = @recipe_group_id
order by r.revision desc;

-- name: CreateRecipe :one
insert into recipes(name, product_id, produced_quantity, created_by_user_id)
values (@name,
        @product_id,
        @produced_quantity,
        @created_by_user_id)
returning recipe_id;

-- name: CreateRecipeIngredients :copyfrom
INSERT INTO recipe_ingredients (recipe_id, product_id, quantity)
VALUES ($1, $2, $3);

-- name: SetCurrentFalse :one
update recipes
set is_current = false
where recipe_id = @recipe_id
returning recipe_id;

-- name: CreateRecipeRevision :one
insert into recipes(name, recipe_group_id, product_id, produced_quantity, revision, created_by_user_id)
values (@name,
        @recipe_group_id,
        @product_id,
        @produced_quantity,
        @revision,
        @created_by_user_id)
returning recipe_id;

-- name: SetRecipeStatusByGroupID :exec
update recipes
set status = @status
where recipe_group_id = @recipe_group_id;
;