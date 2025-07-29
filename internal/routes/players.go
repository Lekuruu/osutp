package routes

import (
	"time"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/xeonx/timeago"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// TODO: Implement persistent state
var printer = message.NewPrinter(language.English)
var lastUpdate = time.Now()
var pageViews = 1

func Players(ctx *common.Context) {
	data := map[string]interface{}{
		"Title":       "osu!DiffCalc - web version",
		"Description": "An attempt to accurately compute beatmap difficulty and player ranking.",
		"LastUpdate":  timeago.English.Format(lastUpdate),
		"PageViews":   printer.Sprintf("%d", pageViews),
		"LoadTime":    0.00069,
	}
	pageViews++
	renderTemplate(ctx.Response, "players", data)
}
