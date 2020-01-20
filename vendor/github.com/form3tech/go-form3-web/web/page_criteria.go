package web

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func BuildPageCriteria(c *gin.Context) PageCriteria {
	pageNumber := c.Query("page[number]")
	pageSize := c.Query("page[size]")

	size, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil || size == 0 {
		size = 1000
	}

	return PageCriteria{
		PageNumber: pageNumber,
		PageSize:   int(size),
	}
}

type PageCriteria struct {
	PageNumber string
	PageSize   int
}

func (p PageCriteria) GetPageNumber(recordCount int) int {
	if strings.ToLower(p.PageNumber) == "last" {
		return lastPageNumber(recordCount, p.PageSize)
	}

	if strings.ToLower(p.PageNumber) == "first" {
		return 0
	}

	page, err := strconv.ParseInt(p.PageNumber, 10, 32)
	if err != nil {
		page = 0
	}

	return int(page)
}
