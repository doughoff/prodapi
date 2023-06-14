-- name: GetRecipes :many
select count(*) over () as full_count,
       r.recipe_id,
       r.recipe_group_id,
       r.status,
       r.name,
       r.revision,
       r.is_current,
       r.created_by_user_id,
       u.name as created_by_user_name,
       r.created_at
from recipes r
    join users u on u.id = r.created_by_user_id
where r.status = any (@status_options::status[])
  and r.name ilike '%' || @search || '%'
order by r.created_at desc
limit @page_limit offset @page_offset;

-- name: GetRecipeIngredients :many
select ri.id,
       ri.recipe_id,
       ri.product_id,
       p.name as product_name,
       ri.quantity
from recipe_ingredients ri
         join products p on ri.product_id = p.id
where recipe_id = any (@recipe_ids::uuid[]);

-- name: GetRecipeByID :one
select r.recipe_id,
       r.recipe_group_id,
       r.name,
       r.status,
       r.revision,
       r.is_current,
       r.created_by_user_id,
       coalesce(u.name, '') as created_by_user_name,
       r.created_at
from recipes r
         left  join users u on u.id = r.created_by_user_id
where r.recipe_id = @recipe_id
limit 1;

-- name: GetRecipesByGroupID :many
select r.recipe_id,
       r.recipe_group_id,
       r.name,
       r.status,
       r.revision,
       r.is_current,
       r.created_by_user_id,
       u.name as created_by_user_name,
       r.created_at
from recipes r
         join users u on u.id = r.created_by_user_id
where r.recipe_group_id = @recipe_group_id
order by r.revision desc;

-- name: CreateRecipe :one
insert into recipes(name, created_by_user_id)
values (@name,
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
insert into recipes(name, recipe_group_id, revision, created_by_user_id)
values (@name,
        @recipe_group_id,
        @revision,
        @created_by_user_id)
returning recipe_id;

-- name: SetRecipeStatusByGroupID :exec
update recipes
set status = @status
where recipe_group_id = @recipe_group_id;
;