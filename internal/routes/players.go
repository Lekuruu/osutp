package routes

import (
	"strings"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
)

const playersPerPage = 50
const countrySliceSize = 20

func Players(ctx *common.Context) {
	pageViews, err := services.IncreasePageViews("players", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	currentPage := GetPageFromQuery(ctx)
	queryOffset := (currentPage - 1) * playersPerPage

	players, err := services.FetchBestPlayers(queryOffset, playersPerPage, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalPlayers, err := services.TotalPlayers(ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}
	totalPages := int(totalPlayers) / playersPerPage
	pagination := NewPaginationData(currentPage, totalPages, playersPerPage, int(totalPlayers))

	bestCountries, err := services.FetchBestCountries(ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	countryPosition := GetCountryPositionFromQuery(ctx)
	offset := countryPosition * countrySliceSize
	end := offset + countrySliceSize
	if end > len(bestCountries) {
		end = len(bestCountries)
	}
	if offset > len(bestCountries) {
		offset = len(bestCountries)
	}
	bestCountrySlice := bestCountries[offset:end]

	data := map[string]interface{}{
		"PageViews":         pageViews,
		"Players":           players,
		"Pagination":        pagination,
		"BestCountries":     bestCountrySlice,
		"CountryPosition":   countryPosition,
		"TotalCountryPages": len(bestCountries) / countrySliceSize,
	}
	renderTemplate(ctx, "players", data)
}

func PlayersByCountry(ctx *common.Context) {
	country := strings.ToUpper(ctx.Vars["country"])
	pageName := "players_" + country

	pageViews, err := services.IncreasePageViews(pageName, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	currentPage := GetPageFromQuery(ctx)
	queryOffset := (currentPage - 1) * playersPerPage

	players, err := services.FetchBestPlayersByCountry(country, queryOffset, playersPerPage, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	totalPlayers, err := services.TotalPlayersByCountry(country, ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}
	totalPages := int(totalPlayers) / playersPerPage
	pagination := NewPaginationData(currentPage, totalPages, playersPerPage, int(totalPlayers))

	bestCountries, err := services.FetchBestCountries(ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	countryPosition := GetCountryPositionFromQuery(ctx)
	offset := countryPosition * countrySliceSize
	end := offset + countrySliceSize
	if end > len(bestCountries) {
		end = len(bestCountries)
	}
	if offset > len(bestCountries) {
		offset = len(bestCountries)
	}
	bestCountrySlice := bestCountries[offset:end]

	data := map[string]interface{}{
		"PageViews":         pageViews,
		"Players":           players,
		"Country":           country,
		"Pagination":        pagination,
		"BestCountries":     bestCountrySlice,
		"CountryPosition":   countryPosition,
		"TotalCountryPages": len(bestCountries) / countrySliceSize,
	}
	renderTemplate(ctx, "country_players", data)
}
