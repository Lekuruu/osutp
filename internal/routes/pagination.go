package routes

import (
	"strconv"

	"github.com/Lekuruu/osutp/internal/common"
)

type PaginationData struct {
	PageViews    int
	CurrentPage  int
	TotalPages   int
	PerPage      int
	TotalResults int
	Start        int
	End          int
}

func NewPaginationData(currentPage, totalPages, perPage, totalResults int) *PaginationData {
	start := (currentPage - 1) * perPage
	end := start + perPage
	if end > totalResults && totalResults != -1 {
		end = totalResults
	}
	if start <= 0 {
		start = 1
	}
	return &PaginationData{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		PerPage:      perPage,
		TotalResults: totalResults,
		Start:        start,
		End:          end,
	}
}

func TestPaginationData() *PaginationData {
	return &PaginationData{
		CurrentPage:  1,
		TotalPages:   100,
		PerPage:      50,
		TotalResults: 0,
		Start:        1,
		End:          50,
	}
}

func GetPageFromQuery(ctx *common.Context) int {
	currentPage := ctx.Request.URL.Query().Get("p")
	if currentPage == "" {
		currentPage = "1"
	}
	currentPageInt, err := strconv.Atoi(currentPage)
	if err != nil {
		currentPageInt = 1
	}
	return currentPageInt
}
