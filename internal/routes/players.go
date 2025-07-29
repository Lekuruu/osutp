package routes

import (
	"time"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/xeonx/timeago"
)

// TODO: Implement "last update"
var lastUpdate = time.Now()

func Players(ctx *common.Context) {
	data := map[string]interface{}{
		"Title":       "osu!DiffCalc - web version",
		"Description": "An attempt to accurately compute beatmap difficulty and player ranking.",
		"LastUpdate":  timeago.English.Format(lastUpdate),
	}
	renderTemplate(ctx.Response, "players", data)
}
