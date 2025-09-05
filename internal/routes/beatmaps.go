package routes

import (
	"strconv"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
	"github.com/Lekuruu/osutp-web/pkg/tp"
)

const defaultMods = tp.NoMod
const beatmapsPerPage = 50

func Beatmaps(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("beatmaps", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	currentMods := GetModsFromQuery(ctx)
	currentPage := GetPageFromQuery(ctx)
	queryOffset := (currentPage - 1) * beatmapsPerPage
	beatmaps, err := services.FetchBeatmapsByDifficulty(queryOffset, beatmapsPerPage, currentMods, ctx.State)
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
	pagination := NewPaginationData(currentPage, totalPages, beatmapsPerPage, int(totalBeatmaps))

	data := map[string]interface{}{
		"PageViews":  pageViews,
		"Pagination": pagination,
		"Beatmaps":   beatmaps,
		"Mods":       currentMods,
	}
	renderTemplate(ctx, "beatmaps", data)
}

func GetModsFromQuery(ctx *common.Context) uint32 {
	mods := defaultMods

	hardRock := ctx.Request.URL.Query().Get("HR")
	easy := ctx.Request.URL.Query().Get("EZ")
	if hardRock == "1" {
		mods |= tp.HardRock
	} else if easy == "1" {
		mods |= tp.Easy
	}

	doubleTime := ctx.Request.URL.Query().Get("DT")
	halfTime := ctx.Request.URL.Query().Get("HT")
	if doubleTime == "1" {
		mods |= tp.DoubleTime
	} else if halfTime == "1" {
		mods |= tp.HalfTime
	}

	return mods
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
