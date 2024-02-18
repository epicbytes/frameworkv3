package http_server

import (
	"github.com/gofiber/fiber/v2"
)

func GetLimitOffset(ctx *fiber.Ctx) (limit, offset int64) {
	return int64(ctx.QueryInt("limit", 0)), int64(ctx.QueryInt("offset", 0))
}
