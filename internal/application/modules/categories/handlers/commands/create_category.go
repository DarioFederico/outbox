package commands

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"outbox/internal/application/commons"
	"outbox/internal/application/modules/categories/dtos"
	"outbox/internal/application/modules/categories/models"
	"outbox/internal/application/modules/categories/models/events"
	"outbox/internal/infrastructure/log"
)

const (
	insertCategory string = `INSERT INTO categories (name, description, created_at, updated_at) VALUE(?,?,?,?)`
	insertMessage  string = `INSERT INTO outbox (type, message, status, created_at, updated_at) VALUE(?,?,?,?,?)`
)

type CreateCategoryCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type createCategoryCommandHandler struct {
	db *sql.DB
}

func NewCreateCategoryCommandHandler(db *sql.DB) commons.RequestHandler {
	return &createCategoryCommandHandler{db: db}
}

func (c *createCategoryCommandHandler) Handle(ctx context.Context, request commons.Request) (commons.Response, error) {
	command := request.(CreateCategoryCommand)
	model := &models.Category{
		Name:        command.Name,
		Description: command.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	//Begin transaction to persist categories and outbox messages
	//if any error was occurred, event never will be sent
	tx, err := c.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	category, err := c.createCategory(ctx, tx, model)
	if err != nil {
		return nil, err
	}

	//Create event and append to outbox message
	event := &events.CreateCategoryEventV1{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Name,
		CreatedAt:   category.CreatedAt,
	}
	err = c.createOutboxMsg(ctx, tx, event)
	if err != nil {
		log.For(ctx).Errorf("error creating outbox message. %+v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.For(ctx).Errorf("error commit transaction. %+v", err)
		return nil, err
	}

	//Map model to dto
	dto := &dtos.CategoryDto{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
	}

	return dto, nil
}

func (c *createCategoryCommandHandler) createCategory(ctx context.Context, tx *sql.Tx, category *models.Category) (*models.Category, error) {
	res, err := tx.ExecContext(ctx, insertCategory,
		category.Name,
		category.Description,
		category.CreatedAt,
		category.UpdatedAt,
	)
	if err != nil {
		log.For(ctx).Errorf("error create category %+v", err)
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.For(ctx).Errorf("error to get insert query result %+v", err)
		return nil, err
	}
	category.ID = lastId
	return category, nil
}

func (c *createCategoryCommandHandler) createOutboxMsg(ctx context.Context, tx *sql.Tx, event events.CategoryEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outbox := models.Outbox{
		Type:      string(event.GetEventName()),
		Message:   string(message),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    models.OutboxStatus_Pending,
	}

	_, err = tx.Exec(insertMessage,
		outbox.Type,
		outbox.Message,
		outbox.Status,
		outbox.CreatedAt,
		outbox.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
