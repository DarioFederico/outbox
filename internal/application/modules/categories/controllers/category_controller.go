package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"outbox/internal/application/commons"
	"outbox/internal/application/modules/categories/handlers/commands"
	"outbox/internal/application/modules/categories/handlers/queries"
	"outbox/internal/application/modules/categories/models"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	mediator commons.Mediator
}

func NewCategoryController(mediator commons.Mediator) *CategoryController {
	return &CategoryController{mediator: mediator}
}

// GetCategoryById
// @Tags Category
// @Summary Get category
// @Description Get category by id
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dtos.CategoryDto
// @Router /categories/{id} [get]
func (cc *CategoryController) GetCategoryById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "BAD REQUEST"})
		return
	}

	query := queries.GetCategoryById{ID: id}
	cc.processHandler(ctx, query)
}

// CreateCategory
// @Tags Category
// @Summary Create category
// @Description Create new category item
// @Accept json
// @Produce json
// @Param CreateCategoryRequestDto body commands.CreateCategoryCommand true "Category data"
// @Success 201 {object} dtos.CategoryDto
// @Router /categories [post]
func (cc *CategoryController) CreateCategory(ctx *gin.Context) {
	var command commands.CreateCategoryCommand
	if err := ctx.Bind(&command); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "BAD REQUEST"})
	} else {
		cc.processHandler(ctx, command)
	}
}

func (cc *CategoryController) processHandler(ctx *gin.Context, request commons.Request) {
	dto, err := cc.mediator.Send(ctx, request)

	if err != nil {
		if errors.Is(err, models.NotFoundError) {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, "ERROR")
		return
	}
	ctx.JSON(http.StatusOK, dto)
}
