// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user_queries.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
insert into users(
                  email, name, password, roles
) values (
          $1,
          $2,
          $3,
          $4
         )
returning id, status, email, name, password, roles, created_at, updated_at
`

type CreateUserParams struct {
	Email    string
	Name     string
	Password string
	Roles    []string
}

func (q *Queries) CreateUser(ctx context.Context, db DBTX, arg *CreateUserParams) (*User, error) {
	row := db.QueryRow(ctx, createUser,
		arg.Email,
		arg.Name,
		arg.Password,
		arg.Roles,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, status, email, name, password, roles, created_at, updated_at from users where email = $1 limit 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, db DBTX, email string) (*User, error) {
	row := db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECt id, status, email, name, password, roles, created_at, updated_at
from users
where id = $1
limit 1
`

func (q *Queries) GetUserByID(ctx context.Context, db DBTX, id pgtype.UUID) (*User, error) {
	row := db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getUsers = `-- name: GetUsers :many
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
where status = any ($1::status[])
and name ilike '%' || $2 || '%'
order by created_at desc
limit $4 offset $3
`

type GetUsersParams struct {
	StatusOptions []Status
	Search        pgtype.Text
	PageOffset    int32
	PageLimit     int32
}

type GetUsersRow struct {
	FullCount int64
	ID        pgtype.UUID
	Status    Status
	Roles     []string
	Name      string
	Email     string
	Password  string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) GetUsers(ctx context.Context, db DBTX, arg *GetUsersParams) ([]*GetUsersRow, error) {
	rows, err := db.Query(ctx, getUsers,
		arg.StatusOptions,
		arg.Search,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetUsersRow{}
	for rows.Next() {
		var i GetUsersRow
		if err := rows.Scan(
			&i.FullCount,
			&i.ID,
			&i.Status,
			&i.Roles,
			&i.Name,
			&i.Email,
			&i.Password,
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

const updateUserByID = `-- name: UpdateUserByID :one
update users set
                 status = $1,
                 email = $2,
                 password= $3,
                 roles = $4,
                 name = $5
where id = $6
returning id, status, email, name, password, roles, created_at, updated_at
`

type UpdateUserByIDParams struct {
	Status   Status
	Email    string
	Password string
	Roles    []string
	Name     string
	ID       pgtype.UUID
}

func (q *Queries) UpdateUserByID(ctx context.Context, db DBTX, arg *UpdateUserByIDParams) (*User, error) {
	row := db.QueryRow(ctx, updateUserByID,
		arg.Status,
		arg.Email,
		arg.Password,
		arg.Roles,
		arg.Name,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}