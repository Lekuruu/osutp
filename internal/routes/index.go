package routes

import (
	"net/http"

	"github.com/Lekuruu/osutp-web/internal/common"
)

func Index(ctx *common.Context) {
	// Redirect to /players page by default
	http.Redirect(ctx.Response, ctx.Request, "/players", http.StatusFound)
}
