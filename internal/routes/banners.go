package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
	"gorm.io/gorm"
)

func Banners(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("banners", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	if playerName := ctx.Request.URL.Query().Get("pn"); playerName != "" {
		player, err := services.FetchPlayerByName(playerName, ctx.State)
		if err != nil && err != gorm.ErrRecordNotFound {
			ctx.Response.WriteHeader(500)
			return
		}

		if player != nil {
			http.Redirect(ctx.Response, ctx.Request, fmt.Sprintf("/banners/%d", player.ID), http.StatusFound)
			return
		}
	}

	if playerId, ok := ctx.Vars["pid"]; ok {
		playerIdInt, err := strconv.Atoi(playerId)
		if err != nil {
			ctx.Response.WriteHeader(400)
			return
		}

		player, err := services.FetchPlayerById(playerIdInt, ctx.State)
		if err != nil && err != gorm.ErrRecordNotFound {
			ctx.Response.WriteHeader(500)
			return
		}

		data := map[string]interface{}{
			"PageViews": pageViews,
			"Player":    player,
		}
		renderTemplate(ctx, "banners_player", data)
		return
	}

	data := map[string]interface{}{"PageViews": pageViews}
	renderTemplate(ctx, "banners", data)
}
