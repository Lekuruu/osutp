package routes

import (
	"html/template"
	"net/url"
	"strconv"

	"github.com/Lekuruu/osutp/internal/common"
)

type PaginationData struct {
	PreviousQuery template.URL
	PageViews     int
	CurrentPage   int
	TotalPages    int
	PerPage       int
	TotalResults  int
	Start         int
	End           int
}

func NewPaginationData(currentPage, totalPages, perPage, totalResults int, query url.Values) *PaginationData {
	start := (currentPage - 1) * perPage
	end := start + perPage
	if end > totalResults && totalResults != -1 {
		end = totalResults
	}
	if start <= 0 {
		start = 1
	}
	return &PaginationData{
		PreviousQuery: buildPreviousQuery(query),
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		PerPage:       perPage,
		TotalResults:  totalResults,
		Start:         start,
		End:           end,
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

func GetCountryPositionFromQuery(ctx *common.Context) int {
	position := ctx.Request.URL.Query().Get("cp")
	if position == "" {
		position = "0"
	}
	positionInt, err := strconv.Atoi(position)
	if err != nil {
		positionInt = 0
	}
	return positionInt
}

func buildPreviousQuery(query url.Values) template.URL {
	if len(query) == 0 {
		return ""
	}
	query.Del("p")
	return template.URL(query.Encode() + "&")
}
