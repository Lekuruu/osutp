package routes

import (
	"strconv"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
)

const scoresPerPage = 50

func Scores(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("scores", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	currentPage := GetPageFromQuery(ctx)
	queryOffset := (currentPage - 1) * scoresPerPage

	bestScores, err := services.FetchBestScores(
		queryOffset, scoresPerPage,
		GetSortColumnFromQuery(ctx), ctx.State,
	)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalScores, err := services.FetchTotalScores(ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalPages := int(totalScores) / scoresPerPage
	pagination := NewPaginationData(
		currentPage, totalPages,
		scoresPerPage, int(totalScores),
	)

	data := map[string]interface{}{
		"PageViews":  pageViews,
		"BestScores": bestScores,
		"Pagination": pagination,
	}
	renderTemplate(ctx, "scores", data)
}

func GetSortColumnFromQuery(ctx *common.Context) string {
	sort := ctx.Request.URL.Query().Get("s")
	sortColumn, err := strconv.Atoi(sort)
	if err != nil {
		return "total_tp DESC"
	}

	switch sortColumn {
	case 0:
		return "total_tp DESC"
	case 1:
		return "aim_tp DESC"
	case 2:
		return "speed_tp DESC"
	case 3:
		return "acc_tp DESC"
	default:
		return "total_tp DESC"
	}
}
