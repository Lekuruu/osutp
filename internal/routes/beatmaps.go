package routes

import (
	"strconv"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
)

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
	pagination := NewPaginationData(currentPageInt, 100, 50, 0)

	data := map[string]interface{}{
		"PageViews":  pageViews,
		"Pagination": pagination,
	}
	renderTemplate(ctx, "beatmaps", data)
}
