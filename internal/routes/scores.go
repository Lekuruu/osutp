package routes

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
)

func Scores(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("scores", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	data := map[string]interface{}{"PageViews": pageViews}
	renderTemplate(ctx, "scores", data)
}
