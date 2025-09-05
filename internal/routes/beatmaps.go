package routes

import (
	"fmt"
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

	filters := []string{}
	filters = ApplyRankedFilter(filters, ctx)
	filters = ApplyRatioFilter(filters, currentMods, ctx)
	filters = ApplyDifficultyFilter("CircleSize", "cs", filters, currentMods, ctx)
	filters = ApplyDifficultyFilter("ApproachRate", "ar", filters, currentMods, ctx)
	filters = ApplyDifficultyFilter("OverallDifficulty", "od", filters, currentMods, ctx)

	beatmaps, err := services.FetchBeatmapsByDifficulty(
		queryOffset, beatmapsPerPage,
		currentMods, filters, ctx.State,
	)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalBeatmaps, err := services.FetchTotalBeatmaps(filters, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}
	totalPages := int(totalBeatmaps) / beatmapsPerPage
	pagination := NewPaginationData(currentPage, totalPages, beatmapsPerPage, int(totalBeatmaps))

	data := map[string]interface{}{
		"Pagination": pagination,
		"PageViews":  pageViews,
		"Beatmaps":   beatmaps,
		"Mods":       currentMods,
		"Query":      ctx.Request.URL.Query(),
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

func ApplyRankedFilter(filters []string, ctx *common.Context) []string {
	statusQuery := ctx.Request.URL.Query().Get("u")
	switch statusQuery {
	case "1":
		filters = append(filters, "status > 0")
	case "2":
		filters = append(filters, "status < 1")
	}
	return filters
}

func ApplyDifficultyFilter(name string, short string, filters []string, mods uint32, ctx *common.Context) []string {
	minDiff := ctx.Request.URL.Query().Get(fmt.Sprintf("%sl", short))
	maxDiff := ctx.Request.URL.Query().Get(fmt.Sprintf("%sh", short))

	// Ensure valid float values
	_, errMin := strconv.ParseFloat(minDiff, 32)
	_, errMax := strconv.ParseFloat(maxDiff, 32)

	if minDiff != "" && errMin == nil {
		filters = append(filters, getDifficultyAttribute(name, mods)+" >= "+minDiff)
	}
	if maxDiff != "" && errMax == nil {
		filters = append(filters, getDifficultyAttribute(name, mods)+" <= "+maxDiff)
	}
	return filters
}

func ApplyRatioFilter(filters []string, mods uint32, ctx *common.Context) []string {
	speed := ctx.Request.URL.Query().Get("s")
	if speed == "" || speed == "50" {
		return filters
	}

	speedValue, err := strconv.ParseFloat(speed, 32)
	if err != nil || speedValue <= 0 || speedValue >= 100 {
		return filters
	}

	// Compute requested ratio
	aimValue := 100 - speedValue
	ratio := speedValue / aimValue

	// Allow a tolerance so it's not too strict
	tolerance := 0.15
	minRatio := ratio - tolerance
	maxRatio := ratio + tolerance

	ratioExpression := fmt.Sprintf(
		"%s / NULLIF(%s, 0)",
		getDifficultyAttribute("SpeedStars", mods),
		getDifficultyAttribute("AimStars", mods),
	)

	filters = append(filters, fmt.Sprintf("%s BETWEEN %f AND %f", ratioExpression, minRatio, maxRatio))
	return filters
}

func getDifficultyAttribute(name string, mods uint32) string {
	return fmt.Sprintf("json_extract(difficulty_attributes, '$.%d.%s')", mods, name)
}
