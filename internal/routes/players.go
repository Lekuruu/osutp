package routes

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
	"github.com/xeonx/timeago"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

func Players(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("players", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	lastUpdate, err := services.PageLastUpdated("players", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	data := map[string]interface{}{
		"LastUpdate": timeago.English.Format(lastUpdate),
		"PageViews":  pageViews,
	}
	renderTemplate(ctx, "players", data)
}
