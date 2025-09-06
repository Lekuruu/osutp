package routes

import (
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

	bestScores, err := services.FetchBestScores(queryOffset, scoresPerPage, "total_tp DESC", ctx.State)
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
	pagination := NewPaginationData(currentPage, totalPages, scoresPerPage, int(totalScores))

	data := map[string]interface{}{
		"PageViews":  pageViews,
		"BestScores": bestScores,
		"Pagination": pagination,
	}
	renderTemplate(ctx, "scores", data)
}
