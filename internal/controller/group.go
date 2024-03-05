package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lmnq/test-thai/internal/service"
	"github.com/lmnq/test-thai/logger"
)

type groupController struct {
	s service.Group
	l logger.Logger
}

func newGroupController(router fiber.Router, l logger.Logger, groupService service.Group) {
	c := &groupController{
		s: groupService,
		l: l,
	}

	r := router.Group("/group")

	r.Post("/", c.create)
	r.Get("/:id", c.get)
	r.Get("/", c.getAll)
	r.Put("/:id", c.update)
	r.Delete("/:id", c.delete)
}

type groupCreateRequest struct {
	GroupName string `json:"group_name"`
}

func (c *groupController) create(ctx *fiber.Ctx) error {
	var req groupCreateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	id, myerr := c.s.Create(ctx.Context(), req.GroupName)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "create group error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (c *groupController) get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get group id param error")
		return errorResponse(ctx, 400, "get group id param error")
	}

	group, myerr := c.s.Get(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get group error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(group)
}

func (c *groupController) getAll(ctx *fiber.Ctx) error {
	groups, myerr := c.s.GetAll(ctx.Context())
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get all groups error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"groups": groups,
	})
}

type groupUpdateRequest struct {
	GroupName string `json:"group_name"`
}

func (c *groupController) update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get group id param error")
		return errorResponse(ctx, 400, "get group id param error")
	}

	var req groupUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	myerr := c.s.Update(ctx.Context(), id, req.GroupName)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "update group error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *groupController) delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get group id param error")
		return errorResponse(ctx, 400, "get group id param error")
	}

	myerr := c.s.Delete(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "delete group error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
