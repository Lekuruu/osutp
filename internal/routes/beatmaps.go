package routes

import (
	"strconv"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
)

const beatmapsPerPage = 50
const defaultMods = 0

func Beatmaps(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("beatmaps", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	currentPage := ctx.Request.URL.Query().Get("p")
	if currentPage == "" {
		currentPage = "1"
	}
	currentPageInt, err := strconv.Atoi(currentPage)
	if err != nil {
		currentPageInt = 1
	}

	queryOffset := (currentPageInt - 1) * beatmapsPerPage
	beatmaps, err := services.FetchBeatmapsByDifficulty(queryOffset, beatmapsPerPage, defaultMods, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalBeatmaps, err := services.FetchTotalBeatmaps(ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}
	totalPages := int(totalBeatmaps) / beatmapsPerPage
	pagination := NewPaginationData(currentPageInt, totalPages, beatmapsPerPage, int(totalBeatmaps))

	data := map[string]interface{}{
		"PageViews":  pageViews,
		"Pagination": pagination,
		"Beatmaps":   beatmaps,
	}
	renderTemplate(ctx, "beatmaps", data)
}
