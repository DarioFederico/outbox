package queries

import (
	"context"
	"database/sql"
	"errors"

	"outbox/internal/application/commons"
	"outbox/internal/application/modules/categories/dtos"
	"outbox/internal/application/modules/categories/models"
	"outbox/internal/infrastructure/log"
)

const (
	getCategoryById string = `SELECT id, name, description, created_at FROM categories WHERE Id = ?`
)

type GetCategoryById struct {
	ID int64
}

type getCategoryByIdHandler struct {
	db *sql.DB
}

func NewGetCategoryByIdHandler(db *sql.DB) commons.RequestHandler {
	return &getCategoryByIdHandler{db: db}
}

func (h *getCategoryByIdHandler) Handle(ctx context.Context, request commons.Request) (commons.Response, error) {
	var dto dtos.CategoryDto
	query, ok := request.(GetCategoryById)
	if !ok {
		err := errors.New("invalid request type for GetCategoryByIdHandler")
		log.For(ctx).Errorf("[GetCategoryByIdHandler:Handle] error to process handler. %+v", err)
		return dto, err
	}

	category, err := h.getCategoryFromDb(ctx, query.ID)
	if err != nil {
		log.For(ctx).Errorf("[GetCategoryByIdHandler:Handle] error to process handler. %+v", err)
		return dto, err
	}

	dto = dtos.CategoryDto{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
	}

	return dto, nil
}

func (h *getCategoryByIdHandler) getCategoryFromDb(ctx context.Context, ID int64) (*models.Category, error) {
	var category models.Category
	err := h.db.QueryRowContext(ctx, getCategoryById, ID).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.NotFoundError
		}
		log.For(ctx).Errorf("error get category %d. %+v", ID, err)
		return nil, err
	}
	return &category, nil

}
