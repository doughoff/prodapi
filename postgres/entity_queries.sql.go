// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: entity_queries.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createEntity = `-- name: CreateEntity :one
INSERT INTO "entities" (name,
                        ruc,
                        ci)
VALUES ($1,
        $2,
        $3)
RETURNING id, status, name, ci, ruc, created_at, updated_at
`

type CreateEntityParams struct {
	Name string
	Ruc  pgtype.Text
	Ci   pgtype.Text
}

func (q *Queries) CreateEntity(ctx context.Context, db DBTX, arg *CreateEntityParams) (*Entity, error) {
	row := db.QueryRow(ctx, createEntity, arg.Name, arg.Ruc, arg.Ci)
	var i Entity
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Ci,
		&i.Ruc,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getEntities = `-- name: GetEntities :many
SELECT COUNT(*) OVER () AS full_count,
       id,
       status,
       name,
       ruc,
       ci,
       created_at,
       updated_at
FROM "entities"
WHERE status = ANY ($1::status[])
   AND (
            name ILIKE '%' || $2 || '%'
        OR ruc ILIKE '%' || $2 || '%'
        OR ci ILIKE '%' || $2 || '%'
    )
ORDER BY created_at DESC
LIMIT $4 OFFSET $3
`

type GetEntitiesParams struct {
	StatusOptions []Status
	Search        pgtype.Text
	PageOffset    int32
	PageLimit     int32
}

type GetEntitiesRow struct {
	FullCount int64
	ID        pgtype.UUID
	Status    Status
	Name      string
	Ruc       pgtype.Text
	Ci        pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) GetEntities(ctx context.Context, db DBTX, arg *GetEntitiesParams) ([]*GetEntitiesRow, error) {
	rows, err := db.Query(ctx, getEntities,
		arg.StatusOptions,
		arg.Search,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetEntitiesRow{}
	for rows.Next() {
		var i GetEntitiesRow
		if err := rows.Scan(
			&i.FullCount,
			&i.ID,
			&i.Status,
			&i.Name,
			&i.Ruc,
			&i.Ci,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntityByCI = `-- name: GetEntityByCI :one
SELECT
    id,
    status,
    name,
    ruc,
    ci,
    created_at,
    updated_at
FROM "entities"
WHERE ci = $1
limit 1
`

type GetEntityByCIRow struct {
	ID        pgtype.UUID
	Status    Status
	Name      string
	Ruc       pgtype.Text
	Ci        pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) GetEntityByCI(ctx context.Context, db DBTX, ci pgtype.Text) (*GetEntityByCIRow, error) {
	row := db.QueryRow(ctx, getEntityByCI, ci)
	var i GetEntityByCIRow
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Ruc,
		&i.Ci,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getEntityByID = `-- name: GetEntityByID :one
SELECT id, status, name, ci, ruc, created_at, updated_at
FROM "entities"
WHERE id = $1
limit 1
`

func (q *Queries) GetEntityByID(ctx context.Context, db DBTX, entityID pgtype.UUID) (*Entity, error) {
	row := db.QueryRow(ctx, getEntityByID, entityID)
	var i Entity
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Ci,
		&i.Ruc,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getEntityByRUC = `-- name: GetEntityByRUC :one
SELECT
    id,
    status,
    name,
    ruc,
    ci,
    created_at,
    updated_at
FROM "entities"
WHERE ruc = $1
limit 1
`

type GetEntityByRUCRow struct {
	ID        pgtype.UUID
	Status    Status
	Name      string
	Ruc       pgtype.Text
	Ci        pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) GetEntityByRUC(ctx context.Context, db DBTX, ruc pgtype.Text) (*GetEntityByRUCRow, error) {
	row := db.QueryRow(ctx, getEntityByRUC, ruc)
	var i GetEntityByRUCRow
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Ruc,
		&i.Ci,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const updateEntityByID = `-- name: UpdateEntityByID :one
UPDATE "entities" SET
                      status = $1,
                      name = $2,
                      ruc = $3,
                      ci = $4
WHERE id = $5
RETURNING id, status, name, ci, ruc, created_at, updated_at
`

type UpdateEntityByIDParams struct {
	Status Status
	Name   string
	Ruc    pgtype.Text
	Ci     pgtype.Text
	ID     pgtype.UUID
}

func (q *Queries) UpdateEntityByID(ctx context.Context, db DBTX, arg *UpdateEntityByIDParams) (*Entity, error) {
	row := db.QueryRow(ctx, updateEntityByID,
		arg.Status,
		arg.Name,
		arg.Ruc,
		arg.Ci,
		arg.ID,
	)
	var i Entity
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Ci,
		&i.Ruc,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}