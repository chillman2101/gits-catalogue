package query

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Page   int
	Limit  int
	Sort   string
	Order  string
	Search string
}

func Parse(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "id")
	order := strings.ToUpper(c.DefaultQuery("order", "ASC"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	return Params{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
		Search: search,
	}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

func (p Params) OrderClause() string {
	return p.Sort + " " + p.Order
}
