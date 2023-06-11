-- name: GetEntities :many
SELECT COUNT(*) OVER () AS full_count,
       id,
       status,
       name,
       ruc,
       ci,
       created_at,
       updated_at
FROM "entities"
WHERE status = ANY (@status_options::status[])
   OR (
            name ILIKE '%' || @search || '%'
        OR ruc ILIKE '%' || @search || '%'
        OR ci ILIKE '%' || @search || '%'
    )
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: GetEntityByID :one
SELECT id,
       status,
       name,
       ruc,
       ci,
       created_at,
       updated_at
FROM "entities"
WHERE id = @entity_id
limit 1;

-- name: CreateEntity :one
INSERT INTO "entities" (name,
                        ruc,
                        ci)
VALUES (@name,
        @ruc,
        @ci)
RETURNING *;

-- name: UpdateEntityByID :one
UPDATE "entities" SET
                      status = @status,
                      name = @name,
                      ruc = @ruc,
                      ci = @ci
WHERE id = @id
RETURNING *;

-- name: GetEntityByRUC :one
SELECT
    id,
    status,
    name,
    ruc,
    ci,
    created_at,
    updated_at
FROM "entities"
WHERE ruc = @ruc
limit 1;

-- name: GetEntityByCI :one
SELECT
    id,
    status,
    name,
    ruc,
    ci,
    created_at,
    updated_at
FROM "entities"
WHERE ci = @ci
limit 1;
