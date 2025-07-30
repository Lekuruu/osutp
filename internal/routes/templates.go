package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
	"github.com/xeonx/timeago"
)

var templates *template.Template

func renderTemplate(ctx *common.Context, tmpl string, pageData map[string]interface{}) {
	lastUpdate, err := services.PageLastUpdated("players", ctx.State)
	if err != nil {
		ctx.Response.WriteHeader(500)
		return
	}

	data := map[string]interface{}{
		"Title":       "osu!DiffCalc - web version",
		"Description": "An attempt to accurately compute beatmap difficulty and player ranking.",
		"LoadTime":    fmt.Sprintf("%.4f", time.Since(ctx.Start).Seconds()),
		"LastUpdate":  timeago.English.Format(lastUpdate),
	}
	for k, v := range pageData {
		data[k] = v
	}

	err = templates.ExecuteTemplate(ctx.Response, tmpl, data)
	if err != nil {
		http.Error(ctx.Response, "Template execution error", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func init() {
	var err error
	templates, err = template.ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}
