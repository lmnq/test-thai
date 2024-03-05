package controller

import "github.com/gofiber/fiber/v2"

type response struct {
	Error string `json:"error"`
}

func errorResponse(ctx *fiber.Ctx, code int, message string) error {
	return ctx.Status(code).JSON(response{message})
}
