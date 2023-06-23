// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: production_order_queries.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrderCycleMovement = `-- name: CreateOrderCycleMovement :one
insert into order_cycles_movements(cycle_id, movement_id)
VALUES ($1, $2)
returning id
`

type CreateOrderCycleMovementParams struct {
	CycleID    pgtype.UUID
	MovementID pgtype.UUID
}

func (q *Queries) CreateOrderCycleMovement(ctx context.Context, db DBTX, arg *CreateOrderCycleMovementParams) (pgtype.UUID, error) {
	row := db.QueryRow(ctx, createOrderCycleMovement, arg.CycleID, arg.MovementID)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const createProductionOrder = `-- name: CreateProductionOrder :one
insert into production_orders(cycles, recipe_id, created_by_user_id)
values ($1, $2, $3)
returning id
`

type CreateProductionOrderParams struct {
	Cycles          int64
	RecipeID        pgtype.UUID
	CreatedByUserID pgtype.UUID
}

func (q *Queries) CreateProductionOrder(ctx context.Context, db DBTX, arg *CreateProductionOrderParams) (pgtype.UUID, error) {
	row := db.QueryRow(ctx, createProductionOrder, arg.Cycles, arg.RecipeID, arg.CreatedByUserID)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

type CreateProductionOrderCyclesParams struct {
	Factor            int64
	ProductionOrderID pgtype.UUID
	CycleOrder        int64
}

const getProductionOrderByID = `-- name: GetProductionOrderByID :one
SELECT po.id,
       po.status,
       po.production_step,
       po.code,
       po.cycles,
       po.output,
       po.recipe_id,
       r.name  as recipe_name,
       r.produced_quantity,
       p.name  as product_name,
       p.unit  as product_unit,
       po.created_by_user_id,
       u.name  as create_by_user_name,
       po.cancelled_by_user_id,
       ud.name as cancelled_by_user_name,
       po.created_at,
       po.updated_at
from production_orders po
         left join recipes r on po.recipe_id = r.recipe_id
         left join products p on r.product_id = p.id
         left join users u on po.created_by_user_id = u.id
         left join users ud on po.cancelled_by_user_id = ud.id
where po.id = $1
`

type GetProductionOrderByIDRow struct {
	ID                  pgtype.UUID
	Status              Status
	ProductionStep      ProductionStep
	Code                pgtype.Text
	Cycles              int64
	Output              pgtype.Int8
	RecipeID            pgtype.UUID
	RecipeName          pgtype.Text
	ProducedQuantity    pgtype.Int8
	ProductName         pgtype.Text
	ProductUnit         NullUnit
	CreatedByUserID     pgtype.UUID
	CreateByUserName    pgtype.Text
	CancelledByUserID   pgtype.UUID
	CancelledByUserName pgtype.Text
	CreatedAt           pgtype.Timestamp
	UpdatedAt           pgtype.Timestamp
}

func (q *Queries) GetProductionOrderByID(ctx context.Context, db DBTX, productionOrderID pgtype.UUID) (*GetProductionOrderByIDRow, error) {
	row := db.QueryRow(ctx, getProductionOrderByID, productionOrderID)
	var i GetProductionOrderByIDRow
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.ProductionStep,
		&i.Code,
		&i.Cycles,
		&i.Output,
		&i.RecipeID,
		&i.RecipeName,
		&i.ProducedQuantity,
		&i.ProductName,
		&i.ProductUnit,
		&i.CreatedByUserID,
		&i.CreateByUserName,
		&i.CancelledByUserID,
		&i.CancelledByUserName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getProductionOrderCycleByID = `-- name: GetProductionOrderCycleByID :one
select id, factor, production_order_id, production_step, cycle_order, completed_at
from production_order_cycles poc
where poc.id = $1
`

func (q *Queries) GetProductionOrderCycleByID(ctx context.Context, db DBTX, cycleID pgtype.UUID) (*ProductionOrderCycle, error) {
	row := db.QueryRow(ctx, getProductionOrderCycleByID, cycleID)
	var i ProductionOrderCycle
	err := row.Scan(
		&i.ID,
		&i.Factor,
		&i.ProductionOrderID,
		&i.ProductionStep,
		&i.CycleOrder,
		&i.CompletedAt,
	)
	return &i, err
}

const getProductionOrderCycles = `-- name: GetProductionOrderCycles :many
select poc.id,
       factor,
       production_order_id,
       production_step,
       cycle_order,
       completed_at
from production_order_cycles poc
where poc.production_order_id = any ($1::uuid[])
`

func (q *Queries) GetProductionOrderCycles(ctx context.Context, db DBTX, productionOrderIds []pgtype.UUID) ([]*ProductionOrderCycle, error) {
	rows, err := db.Query(ctx, getProductionOrderCycles, productionOrderIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ProductionOrderCycle{}
	for rows.Next() {
		var i ProductionOrderCycle
		if err := rows.Scan(
			&i.ID,
			&i.Factor,
			&i.ProductionOrderID,
			&i.ProductionStep,
			&i.CycleOrder,
			&i.CompletedAt,
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

const getProductionOrderMovements = `-- name: GetProductionOrderMovements :many
select ocm.cycle_id    as production_order_cycle_id,
       ocm.movement_id as stock_movement_id
from production_order_cycles poc
         join order_cycles_movements ocm on poc.id = ocm.cycle_id
where poc.production_order_id = $1
`

type GetProductionOrderMovementsRow struct {
	ProductionOrderCycleID pgtype.UUID
	StockMovementID        pgtype.UUID
}

func (q *Queries) GetProductionOrderMovements(ctx context.Context, db DBTX, productionOrderID pgtype.UUID) ([]*GetProductionOrderMovementsRow, error) {
	rows, err := db.Query(ctx, getProductionOrderMovements, productionOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetProductionOrderMovementsRow{}
	for rows.Next() {
		var i GetProductionOrderMovementsRow
		if err := rows.Scan(&i.ProductionOrderCycleID, &i.StockMovementID); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductionOrders = `-- name: GetProductionOrders :many
SELECT count(*) over ()    as full_count,
       po.id,
       po.status,
       po.production_step,
       po.code,
       po.cycles,
       po.output,
       po.recipe_id,
       r.name              as recipe_name,
       r.produced_quantity as produced_quantity,
       p.name              as product_name,
       p.unit              as product_unit,
       po.created_by_user_id,
       u.name              as create_by_user_name,
       po.cancelled_by_user_id,
       ud.name             as cancelled_by_user_name,
       po.created_at,
       po.updated_at
from production_orders po
         left join recipes r on po.recipe_id = r.recipe_id
         left join products p on r.product_id = p.id
         left join users u on po.created_by_user_id = u.id
         left join users ud on po.cancelled_by_user_id = ud.id
where po.status = any ($1::status[])
  and (
            r.name ilike '%' || $2 || '%'
        or p.name ilike '%' || $2 || '%'
    )
  and po.created_at >= $3
order by po.created_at desc
limit $5 offset $4
`

type GetProductionOrdersParams struct {
	StatusOptions []Status
	Search        pgtype.Text
	StartDate     pgtype.Timestamp
	PageOffset    int32
	PageLimit     int32
}

type GetProductionOrdersRow struct {
	FullCount           int64
	ID                  pgtype.UUID
	Status              Status
	ProductionStep      ProductionStep
	Code                pgtype.Text
	Cycles              int64
	Output              pgtype.Int8
	RecipeID            pgtype.UUID
	RecipeName          pgtype.Text
	ProducedQuantity    pgtype.Int8
	ProductName         pgtype.Text
	ProductUnit         NullUnit
	CreatedByUserID     pgtype.UUID
	CreateByUserName    pgtype.Text
	CancelledByUserID   pgtype.UUID
	CancelledByUserName pgtype.Text
	CreatedAt           pgtype.Timestamp
	UpdatedAt           pgtype.Timestamp
}

func (q *Queries) GetProductionOrders(ctx context.Context, db DBTX, arg *GetProductionOrdersParams) ([]*GetProductionOrdersRow, error) {
	rows, err := db.Query(ctx, getProductionOrders,
		arg.StatusOptions,
		arg.Search,
		arg.StartDate,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetProductionOrdersRow{}
	for rows.Next() {
		var i GetProductionOrdersRow
		if err := rows.Scan(
			&i.FullCount,
			&i.ID,
			&i.Status,
			&i.ProductionStep,
			&i.Code,
			&i.Cycles,
			&i.Output,
			&i.RecipeID,
			&i.RecipeName,
			&i.ProducedQuantity,
			&i.ProductName,
			&i.ProductUnit,
			&i.CreatedByUserID,
			&i.CreateByUserName,
			&i.CancelledByUserID,
			&i.CancelledByUserName,
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

const updateProductionOrder = `-- name: UpdateProductionOrder :exec
update production_orders
set status          = $1,
    production_step = $2
where id = $3
`

type UpdateProductionOrderParams struct {
	Status         Status
	ProductionStep ProductionStep
	ID             pgtype.UUID
}

func (q *Queries) UpdateProductionOrder(ctx context.Context, db DBTX, arg *UpdateProductionOrderParams) error {
	_, err := db.Exec(ctx, updateProductionOrder, arg.Status, arg.ProductionStep, arg.ID)
	return err
}

const updateProductionOrderCycle = `-- name: UpdateProductionOrderCycle :exec
update production_order_cycles
set  production_step = $1,
     completed_at = $2
where id = $3
`

type UpdateProductionOrderCycleParams struct {
	ProductionStep ProductionStep
	CompletedAt    pgtype.Timestamp
	ID             pgtype.UUID
}

func (q *Queries) UpdateProductionOrderCycle(ctx context.Context, db DBTX, arg *UpdateProductionOrderCycleParams) error {
	_, err := db.Exec(ctx, updateProductionOrderCycle, arg.ProductionStep, arg.CompletedAt, arg.ID)
	return err
}
