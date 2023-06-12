package routes

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodapi/postgres"
	"github.com/hoffax/prodapi/server/dto"
	"github.com/hoffax/prodapi/server/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type getAllProductsQuery struct {
	StatusOptions []string `query:"status"`
	Search        string   `query:"search"`
	Limit         int32    `query:"limit"`
	Offset        int32    `query:"offset"`
}

func (r *RouteManager) getAllProducts(c *fiber.Ctx, tx *pgx.Tx) error {
	params := new(getAllProductsQuery)
	if err := c.QueryParser(params); err != nil {
		return types.NewInvalidParamsError("invalid query params")
	}

	statusOptions := make([]postgres.Status, len(params.StatusOptions))
	for i, status := range params.StatusOptions {
		statusOptions[i] = postgres.Status(status)
		if !statusOptions[i].Valid() {
			return types.NewInvalidParamsError("invalid status option")
		}
	}
	products, err := r.db.GetProducts(c.Context(), *tx, &postgres.GetProductsParams{
		PageOffset:    params.Offset,
		PageLimit:     params.Limit,
		Search:        pgtype.Text(sql.NullString{String: params.Search, Valid: true}),
		StatusOptions: statusOptions,
	})
	if err != nil {
		return err
	}

	var totalCount int64
	if len(products) > 0 {
		totalCount = products[0].FullCount
	}
	resultRows := make([]*dto.ProductDTO, len(products))
	for i := range resultRows {
		resultRows[i] = &dto.ProductDTO{
			ID:               products[i].ID,
			Status:           products[i].Status,
			Name:             products[i].Name,
			Unit:             products[i].Unit,
			Barcode:          products[i].Barcode,
			ConversionFactor: products[i].ConversionFactor,
			BatchControl:     products[i].BatchControl,
			CreatedAt:        products[i].CreatedAt.Time,
			UpdatedAt:        products[i].UpdatedAt.Time,
		}
	}

	return c.JSON(struct {
		TotalCount int64             `json:"totalCount"`
		Items      []*dto.ProductDTO `json:"items"`
	}{
		TotalCount: totalCount,
		Items:      resultRows,
	})
}

func (r *RouteManager) getProductByID(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	productID := pgtype.UUID{}
	if err := productID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	product, err := r.db.GetProductByID(c.Context(), *tx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.NewNotFoundError()
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToProductDTO(product))
}

func (r *RouteManager) getProductByIBarcode(c *fiber.Ctx, tx *pgx.Tx) error {
	barcode := c.Params("id")

	product, err := r.db.GetProductByBarcode(c.Context(), *tx, barcode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.NewNotFoundError()
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToProductDTO(product))
}

type createProductBody struct {
	Name             string        `json:"name" validate:"required,min=1,max=255"`
	Unit             postgres.Unit `json:"unit" validate:"required,min=1,max=255"`
	Barcode          string        `json:"barcode" validate:"required,min=1,max=100"`
	ConversionFactor int32         `json:"conversionFactor" validate:"required"`
	BatchControl     bool          `json:"batchControl"`
}

func (r *RouteManager) createProduct(c *fiber.Ctx, tx *pgx.Tx) error {
	body := new(createProductBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	if !body.Unit.Valid() {
		return types.NewInvalidParamsError("invalid value for status")
	}

	_, err := r.db.GetProductByBarcode(c.Context(), *tx, body.Barcode)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	} else {
		return types.NewInvalidParamsError("product with this barcode already exists")
	}

	product, err := r.db.CreateProduct(c.Context(), *tx, &postgres.CreateProductParams{
		Name:             body.Name,
		Unit:             body.Unit,
		Barcode:          body.Barcode,
		ConversionFactor: body.ConversionFactor,
		BatchControl:     body.BatchControl,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToProductDTO(product))
}

type updateProductBody struct {
	Status           postgres.Status `json:"status" validate:"required"`
	Name             string          `json:"name" validate:"required,min=1,max=255"`
	Unit             postgres.Unit   `json:"unit" validate:"required,min=1,max=255"`
	Barcode          string          `json:"barcode" validate:"required,min=1,max=100"`
	ConversionFactor int32           `json:"conversionFactor" validate:"required"`
	BatchControl     bool            `json:"batchControl"`
}

func (r *RouteManager) updateProduct(c *fiber.Ctx, tx *pgx.Tx) error {
	idParam := c.Params("id")
	productID := pgtype.UUID{}
	if err := productID.Scan(idParam); err != nil {
		return types.NewInvalidParamsError("invalid uuid on id url param")
	}

	body := new(updateProductBody)
	if err := c.BodyParser(body); err != nil {
		return types.NewInvalidBodyError()
	}
	if err := r.validate.Struct(body); err != nil {
		return err
	}

	if !body.Status.Valid() {
		return types.NewInvalidParamsError("invalid value for status")
	}
	if !body.Unit.Valid() {
		return types.NewInvalidParamsError("invalid value for status")
	}

	product, err := r.db.GetProductByBarcode(c.Context(), *tx, body.Barcode)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	} else {
		if !bytes.Equal(product.ID.Bytes[:], productID.Bytes[:]) {
			return types.NewInvalidParamsError("product with this barcode already exists")
		}
	}

	product, err = r.db.UpdateProductByID(c.Context(), *tx, &postgres.UpdateProductByIDParams{
		ID:               productID,
		Status:           body.Status,
		Name:             body.Name,
		Unit:             body.Unit,
		Barcode:          body.Barcode,
		ConversionFactor: body.ConversionFactor,
		BatchControl:     body.BatchControl,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToProductDTO(product))
}

func (r *RouteManager) RegisterProductRoutes() {
	r.app.Get("/products", r.dbWrapper.WithTransaction(r.getAllProducts))
	r.app.Get("/products/:id", r.dbWrapper.WithTransaction(r.getProductByID))
	r.app.Post("/products", r.dbWrapper.WithTransaction(r.createProduct))
	r.app.Put("/products/:id", r.dbWrapper.WithTransaction(r.updateProduct))
	r.app.Get("/check_barcode/:id", r.dbWrapper.WithTransaction(r.getProductByIBarcode))
}
