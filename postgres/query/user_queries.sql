-- name: GetUsers :many
select
    count(*) over () as full_count,
    id,
    status,
    roles,
    name,
    email,
    password,
    created_at,
    updated_at
from users
where status = any (@status_options::status[])
and name ilike '%' || @search || '%'
order by created_at desc
limit @page_limit offset @page_offset;

-- name: GetUserByID :one
SELECt *
from users
where id = @id
limit 1;

-- name: CreateUser :one
insert into users(
                  email, name, password, roles
) values (
          @email,
          @name,
          @password,
          @roles
         )
returning *;

-- name: UpdateUserByID :one
update users set
                 status = @status,
                 email = @email,
                 password= @password,
                 roles = @roles,
                 name = @name
where id = @id
returning *;

-- name: GetUserByEmail :one
select * from users where email = @email limit 1;