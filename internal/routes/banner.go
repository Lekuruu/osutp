package routes

import (
	"fmt"
	"image/png"
	"strconv"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
	"github.com/Lekuruu/osutp/pkg/banners"
)

func BannerImage(ctx *common.Context) {
	playerIdString, ok := ctx.Vars["pid"]
	if !ok {
		ctx.Response.WriteHeader(400)
		return
	}

	playerId, err := strconv.Atoi(playerIdString)
	if err != nil {
		ctx.Response.WriteHeader(400)
		return
	}

	player, err := services.FetchPlayerById(playerId, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	bannerPlayer := banners.NewPlayer(
		player.Name,
		player.CountryName(),
		player.Country,
		fmt.Sprintf("%s.gif", player.Country),
		player.GlobalRank,
		player.CountryRank,
		int(player.TotalTp),
	)
	bannerIdQuery := ctx.Request.URL.Query().Get("id")
	bannerId, err := strconv.Atoi(bannerIdQuery)
	if err != nil {
		bannerId = 0
	}

	banner := GetBannerById(bannerId, bannerPlayer)
	img := banner.Render()

	ctx.Response.Header().Set("Content-Type", "image/png")
	err = png.Encode(ctx.Response, img)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}
}

func GetBannerById(id int, player banners.Player) banners.Banner {
	switch id {
	case 0:
		return banners.NewCleanStyleBanner(player)
	case 1:
		return banners.NewCleanStyleCenteredBanner(player)
	case 2:
		return banners.NewCleanStyleOneLineBanner(player)
	default:
		return banners.NewCleanStyleBanner(player)
	}
}

func init() {
	banners.TahomaBoldFontPath = "./web/static/fonts/tahoma-bold.ttf"
	banners.TahomaFontPath = "./web/static/fonts/tahoma.ttf"
}
