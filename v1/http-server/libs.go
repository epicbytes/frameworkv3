package http_server

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetLimitOffset(ctx *fiber.Ctx) (limit, offset int64) {
	var err error
	limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
	if err != nil {
		limit = 0
	}
	offset, err = strconv.ParseInt(ctx.Query("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	return limit, offset
}
