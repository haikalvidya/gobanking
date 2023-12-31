package payload

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetPagination(c echo.Context) *PaginationRequest {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("perpage"))
	return GetPaginationV2(page, perPage)
}

func GetPaginationV2(page int, perPage int) *PaginationRequest {
	if page == 0 || page < 1 {
		page = 1
	}

	if perPage == 0 || perPage < 1 {
		perPage = 10
	}

	var offset int

	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * perPage
	}

	limit := perPage

	return &PaginationRequest{
		Page:    int64(page),
		PerPage: int64(perPage),
		Limit:   int64(limit),
		Offset:  int64(offset),
	}
}

func CalculateMetaPagination(totalData int64, p *PaginationRequest) *Pagination {

	tmp := float64(totalData) / float64(p.PerPage)
	lastPage := int64(math.Ceil(tmp))

	isLoadMore := p.Page != lastPage

	if p.Page > lastPage {
		isLoadMore = false
	}

	return &Pagination{
		PerPage:     p.PerPage,
		CurrentPage: p.Page,
		IsLoadMore:  isLoadMore,
		LastPage:    lastPage,
	}
}
