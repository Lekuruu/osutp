package routes

import (
	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
)

func Changelog(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("changelog", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	changelogs, err := services.FetchChangelogs(100, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	data := map[string]interface{}{
		"PageViews":  pageViews,
		"Changelogs": changelogs,
	}
	renderTemplate(ctx, "changelog", data)
}
