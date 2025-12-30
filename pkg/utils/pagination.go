package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func GetPagination(c *fiber.Ctx) Pagination {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	return Pagination{
		Limit: limit,
		Page:  page,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
