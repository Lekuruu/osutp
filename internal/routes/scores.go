package routes

import (
	"fmt"
	"strconv"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
)

const scoresPerPage = 50

func Scores(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("scores", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	playerName := ctx.Request.URL.Query().Get("pn")
	if playerName != "" {
		ScoresByPlayerName(playerName, ctx)
		return
	}

	playerId := ctx.Request.URL.Query().Get("pid")
	if playerId != "" {
		ScoresByPlayer(pageViews, playerId, ctx)
		return
	}

	beatmapId := ctx.Request.URL.Query().Get("bid")
	if beatmapId != "" {
		ScoresByBeatmap(pageViews, beatmapId, ctx)
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

func ScoresByPlayer(pageViews int64, playerIdQuery string, ctx *common.Context) {
	playerId, err := strconv.Atoi(playerIdQuery)
	if err != nil {
		ctx.Response.WriteHeader(400)
		return
	}

	player, err := services.FetchPlayerById(playerId, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	currentPage := GetPageFromQuery(ctx)
	queryOffset := (currentPage - 1) * scoresPerPage

	bestScores, err := services.FetchRangePersonalBestScores(
		player.ID, queryOffset, scoresPerPage,
		GetSortColumnFromQuery(ctx), ctx.State,
	)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalScores, err := services.FetchTotalPersonalBestScores(player.ID, ctx.State)
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
		"Player":     player,
	}
	renderTemplate(ctx, "scores_player", data)
}

func ScoresByPlayerName(playerName string, ctx *common.Context) {
	player, err := services.FetchPlayerByName(playerName, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	ctx.Response.Header().Set("Location", fmt.Sprintf("/scores?pid=%d", player.ID))
	ctx.Response.WriteHeader(302)
}

func ScoresByBeatmap(pageViews int64, beatmapIdQuery string, ctx *common.Context) {
	beatmapId, err := strconv.Atoi(beatmapIdQuery)
	if err != nil {
		ctx.Response.WriteHeader(400)
		return
	}

	beatmap, err := services.FetchBeatmapById(beatmapId, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	currentPage := GetPageFromQuery(ctx)
	queryOffset := (currentPage - 1) * scoresPerPage

	bestScores, err := services.FetchBestScoresByBeatmap(
		beatmap.ID, queryOffset, scoresPerPage,
		GetSortColumnFromQuery(ctx), ctx.State,
	)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalScores, err := services.FetchTotalScoresByBeatmap(beatmap.ID, ctx.State)
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
		"Beatmap":    beatmap,
	}
	renderTemplate(ctx, "scores_beatmap", data)
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
