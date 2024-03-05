package controller

import (
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lmnq/test-thai/internal/service"
	"github.com/lmnq/test-thai/logger"
)

func New(f *fiber.App, l logger.Logger, services *service.Service) {
	// options
	f.Use(recover.New())
	f.Use(middleware.Logger())

	// router
	router := f.Group("/")

	// init routes
	newItemController(router, l, services.Item)
	newCategoryController(router, l, services.Category)
	newGroupController(router, l, services.Group)
	newItemDetailController(router, l, services.ItemDetail)
}
