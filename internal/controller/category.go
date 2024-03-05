package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lmnq/test-thai/internal/service"
	"github.com/lmnq/test-thai/logger"
)

type categoryController struct {
	s service.Category
	l logger.Logger
}

func newCategoryController(router fiber.Router, l logger.Logger, categoryService service.Category) {
	c := &categoryController{
		s: categoryService,
		l: l,
	}

	r := router.Group("/category")

	r.Post("/", c.create)
	r.Get("/:id", c.get)
	r.Get("/", c.getAll)
	r.Put("/:id", c.update)
	r.Delete("/:id", c.delete)
}

type categoryCreateRequest struct {
	CategoryName string `json:"category_name"`
}

func (c *categoryController) create(ctx *fiber.Ctx) error {
	var req categoryCreateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	id, myerr := c.s.Create(ctx.Context(), req.CategoryName)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "create category error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (c *categoryController) get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get category id param error")
		return errorResponse(ctx, 400, "get category id param error")
	}

	category, myerr := c.s.Get(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get category error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(category)
}

func (c *categoryController) getAll(ctx *fiber.Ctx) error {
	categories, myerr := c.s.GetAll(ctx.Context())
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get all categories error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(categories)
}

type categoryUpdateRequest struct {
	CategoryName string `json:"category_name"`
}

func (c *categoryController) update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get category id param error")
		return errorResponse(ctx, 400, "get category id param error")
	}

	var req categoryUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	myerr := c.s.Update(ctx.Context(), id, req.CategoryName)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "update category error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}

func (c *categoryController) delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get category id param error")
		return errorResponse(ctx, 400, "get category id param error")
	}

	myerr := c.s.Delete(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "delete category error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}
