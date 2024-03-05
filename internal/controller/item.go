package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lmnq/test-thai/internal/service"
	"github.com/lmnq/test-thai/logger"
)

type itemController struct {
	s service.Item
	l logger.Logger
}

func newItemController(router fiber.Router, l logger.Logger, itemService service.Item) {
	c := &itemController{
		s: itemService,
		l: l,
	}

	r := router.Group("/item")

	r.Post("/", c.create)
	r.Get("/:id", c.get)
	r.Get("/", c.getAll)
	r.Put("/:id", c.update)
	r.Delete("/:id", c.delete)
}

type itemCreateRequest struct {
	ItemName string `json:"item_name"`
}

func (c *itemController) create(ctx *fiber.Ctx) error {
	var req itemCreateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	id, myerr := c.s.Create(ctx.Context(), req.ItemName)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "create item error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (c *itemController) get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get item id param error")
		return errorResponse(ctx, 400, "get item id param error")
	}

	item, myerr := c.s.Get(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get item error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(item)
}

func (c *itemController) getAll(ctx *fiber.Ctx) error {
	items, myerr := c.s.GetAll(ctx.Context())
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get all items error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(items)
}

type itemUpdateRequest struct {
	ItemName string `json:"item_name"`
}

func (c *itemController) update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get item id param error")
		return errorResponse(ctx, 400, "get item id param error")
	}

	var req itemUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	myerr := c.s.Update(ctx.Context(), id, req.ItemName)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "update item error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}

func (c *itemController) delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get item id param error")
		return errorResponse(ctx, 400, "get item id param error")
	}

	myerr := c.s.Delete(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "delete item error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}
