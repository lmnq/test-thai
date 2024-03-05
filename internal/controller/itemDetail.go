package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lmnq/test-thai/internal/model"
	"github.com/lmnq/test-thai/internal/service"
	"github.com/lmnq/test-thai/logger"
)

type itemDetailController struct {
	s service.ItemDetail
	l logger.Logger
}

func newItemDetailController(router fiber.Router, l logger.Logger, itemDetailService service.ItemDetail) {
	c := &itemDetailController{
		s: itemDetailService,
		l: l,
	}

	r := router.Group("/item_detail")

	r.Post("/", c.create)
	r.Get("/:id", c.get)
	r.Get("/", c.getAllFilter)
	r.Put("/:id", c.update)
	r.Delete("/:id", c.delete)
}

type itemDetailCreateRequest struct {
	ItemName   string  `json:"item_name"`
	GroupID    int     `json:"group_id"`
	CategoryID int     `json:"category_id"`
	Cost       float64 `json:"cost"`
	Price      float64 `json:"price"`
	Sort       int     `json:"sort"`
}

func (c *itemDetailController) create(ctx *fiber.Ctx) error {
	var req itemDetailCreateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	id, myerr := c.s.Create(ctx.Context(), &model.ItemDetail{
		GroupID:    req.GroupID,
		CategoryID: req.CategoryID,
		Cost:       req.Cost,
		Price:      req.Price,
		Sort:       req.Sort,
	}, req.ItemName)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "create item detail error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func (c *itemDetailController) get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get item detail id param error")
		return errorResponse(ctx, 400, "get item detail id param error")
	}

	itemDetail, myerr := c.s.Get(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get item detail error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(itemDetail)
}

type itemDetailFilterParams struct {
	ID           int    `query:"id"`
	ItemName     string `query:"item_name"`
	CategoryName string `query:"category_name"`
	GroupName    string `query:"group_name"`
}

func (c *itemDetailController) getAllFilter(ctx *fiber.Ctx) error {
	var params itemDetailFilterParams

	if err := ctx.QueryParser(&params); err != nil {
		c.l.Error(err, "query parser error")
		return errorResponse(ctx, 400, "query parser error")
	}

	itemDetails, myerr := c.s.GetAllFilter(ctx.Context(), &model.ItemDetailFilter{
		ID:           &params.ID,
		ItemName:     &params.ItemName,
		CategoryName: &params.CategoryName,
		GroupName:    &params.GroupName,
	})
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "get item detail list error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(itemDetails)
}

type itemDetailUpdateRequest struct {
	ItemName   string  `json:"item_name"`
	GroupID    int     `json:"group_id"`
	CategoryID int     `json:"category_id"`
	Cost       float64 `json:"cost"`
	Price      float64 `json:"price"`
	Sort       int     `json:"sort"`
}

func (c *itemDetailController) update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get item detail id param error")
		return errorResponse(ctx, 400, "get item detail id param error")
	}

	var req itemDetailUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.l.Error(err, "body parser error")
		return errorResponse(ctx, 400, "body parser error")
	}

	myerr := c.s.Update(ctx.Context(), id, req.ItemName, &model.ItemDetail{
		GroupID:    req.GroupID,
		CategoryID: req.CategoryID,
		Cost:       req.Cost,
		Price:      req.Price,
		Sort:       req.Sort,
	})
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "update item detail error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}

func (c *itemDetailController) delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.l.Error(err, "get item detail id param error")
		return errorResponse(ctx, 400, "get item detail id param error")
	}

	myerr := c.s.Delete(ctx.Context(), id)
	if myerr.IsErr() {
		c.l.Error(myerr.Err, "delete item detail error")
		return errorResponse(ctx, myerr.Code, myerr.Message)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}
