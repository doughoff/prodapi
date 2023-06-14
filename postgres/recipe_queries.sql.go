// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: recipe_queries.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRecipe = `-- name: CreateRecipe :one
insert into recipes(name, created_by_user_id)
values ($1,
        $2)
returning recipe_id
`

type CreateRecipeParams struct {
	Name            string
	CreatedByUserID pgtype.UUID
}

func (q *Queries) CreateRecipe(ctx context.Context, db DBTX, arg *CreateRecipeParams) (pgtype.UUID, error) {
	row := db.QueryRow(ctx, createRecipe, arg.Name, arg.CreatedByUserID)
	var recipe_id pgtype.UUID
	err := row.Scan(&recipe_id)
	return recipe_id, err
}

type CreateRecipeIngredientsParams struct {
	RecipeID  pgtype.UUID
	ProductID pgtype.UUID
	Quantity  int32
}

const createRecipeRevision = `-- name: CreateRecipeRevision :one
insert into recipes(name, recipe_group_id, revision, created_by_user_id)
values ($1,
        $2,
        $3,
        $4)
returning recipe_id
`

type CreateRecipeRevisionParams struct {
	Name            string
	RecipeGroupID   pgtype.UUID
	Revision        int32
	CreatedByUserID pgtype.UUID
}

func (q *Queries) CreateRecipeRevision(ctx context.Context, db DBTX, arg *CreateRecipeRevisionParams) (pgtype.UUID, error) {
	row := db.QueryRow(ctx, createRecipeRevision,
		arg.Name,
		arg.RecipeGroupID,
		arg.Revision,
		arg.CreatedByUserID,
	)
	var recipe_id pgtype.UUID
	err := row.Scan(&recipe_id)
	return recipe_id, err
}

const getRecipeByID = `-- name: GetRecipeByID :one
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
where r.recipe_id = $1
limit 1
`

type GetRecipeByIDRow struct {
	RecipeID          pgtype.UUID
	RecipeGroupID     pgtype.UUID
	Name              string
	Status            Status
	Revision          int32
	IsCurrent         bool
	CreatedByUserID   pgtype.UUID
	CreatedByUserName string
	CreatedAt         pgtype.Timestamp
}

func (q *Queries) GetRecipeByID(ctx context.Context, db DBTX, recipeID pgtype.UUID) (*GetRecipeByIDRow, error) {
	row := db.QueryRow(ctx, getRecipeByID, recipeID)
	var i GetRecipeByIDRow
	err := row.Scan(
		&i.RecipeID,
		&i.RecipeGroupID,
		&i.Name,
		&i.Status,
		&i.Revision,
		&i.IsCurrent,
		&i.CreatedByUserID,
		&i.CreatedByUserName,
		&i.CreatedAt,
	)
	return &i, err
}

const getRecipeIngredients = `-- name: GetRecipeIngredients :many
select ri.id,
       ri.recipe_id,
       ri.product_id,
       p.name as product_name,
       ri.quantity
from recipe_ingredients ri
         join products p on ri.product_id = p.id
where recipe_id = any ($1::uuid[])
`

type GetRecipeIngredientsRow struct {
	ID          pgtype.UUID
	RecipeID    pgtype.UUID
	ProductID   pgtype.UUID
	ProductName string
	Quantity    int32
}

func (q *Queries) GetRecipeIngredients(ctx context.Context, db DBTX, recipeIds []pgtype.UUID) ([]*GetRecipeIngredientsRow, error) {
	rows, err := db.Query(ctx, getRecipeIngredients, recipeIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeIngredientsRow{}
	for rows.Next() {
		var i GetRecipeIngredientsRow
		if err := rows.Scan(
			&i.ID,
			&i.RecipeID,
			&i.ProductID,
			&i.ProductName,
			&i.Quantity,
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

const getRecipes = `-- name: GetRecipes :many
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
where r.status = any ($1::status[])
  and r.name ilike '%' || $2 || '%'
order by r.created_at desc
limit $4 offset $3
`

type GetRecipesParams struct {
	StatusOptions []Status
	Search        pgtype.Text
	PageOffset    int32
	PageLimit     int32
}

type GetRecipesRow struct {
	FullCount         int64
	RecipeID          pgtype.UUID
	RecipeGroupID     pgtype.UUID
	Status            Status
	Name              string
	Revision          int32
	IsCurrent         bool
	CreatedByUserID   pgtype.UUID
	CreatedByUserName string
	CreatedAt         pgtype.Timestamp
}

func (q *Queries) GetRecipes(ctx context.Context, db DBTX, arg *GetRecipesParams) ([]*GetRecipesRow, error) {
	rows, err := db.Query(ctx, getRecipes,
		arg.StatusOptions,
		arg.Search,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipesRow{}
	for rows.Next() {
		var i GetRecipesRow
		if err := rows.Scan(
			&i.FullCount,
			&i.RecipeID,
			&i.RecipeGroupID,
			&i.Status,
			&i.Name,
			&i.Revision,
			&i.IsCurrent,
			&i.CreatedByUserID,
			&i.CreatedByUserName,
			&i.CreatedAt,
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

const getRecipesByGroupID = `-- name: GetRecipesByGroupID :many
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
where r.recipe_group_id = $1
order by r.revision desc
`

type GetRecipesByGroupIDRow struct {
	RecipeID          pgtype.UUID
	RecipeGroupID     pgtype.UUID
	Name              string
	Status            Status
	Revision          int32
	IsCurrent         bool
	CreatedByUserID   pgtype.UUID
	CreatedByUserName string
	CreatedAt         pgtype.Timestamp
}

func (q *Queries) GetRecipesByGroupID(ctx context.Context, db DBTX, recipeGroupID pgtype.UUID) ([]*GetRecipesByGroupIDRow, error) {
	rows, err := db.Query(ctx, getRecipesByGroupID, recipeGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipesByGroupIDRow{}
	for rows.Next() {
		var i GetRecipesByGroupIDRow
		if err := rows.Scan(
			&i.RecipeID,
			&i.RecipeGroupID,
			&i.Name,
			&i.Status,
			&i.Revision,
			&i.IsCurrent,
			&i.CreatedByUserID,
			&i.CreatedByUserName,
			&i.CreatedAt,
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

const setCurrentFalse = `-- name: SetCurrentFalse :one
update recipes
set is_current = false
where recipe_id = $1
returning recipe_id
`

func (q *Queries) SetCurrentFalse(ctx context.Context, db DBTX, recipeID pgtype.UUID) (pgtype.UUID, error) {
	row := db.QueryRow(ctx, setCurrentFalse, recipeID)
	var recipe_id pgtype.UUID
	err := row.Scan(&recipe_id)
	return recipe_id, err
}

const setRecipeStatusByGroupID = `-- name: SetRecipeStatusByGroupID :exec
update recipes
set status = $1
where recipe_group_id = $2
`

type SetRecipeStatusByGroupIDParams struct {
	Status        Status
	RecipeGroupID pgtype.UUID
}

func (q *Queries) SetRecipeStatusByGroupID(ctx context.Context, db DBTX, arg *SetRecipeStatusByGroupIDParams) error {
	_, err := db.Exec(ctx, setRecipeStatusByGroupID, arg.Status, arg.RecipeGroupID)
	return err
}
