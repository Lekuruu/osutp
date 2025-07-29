package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Lekuruu/osutp-web/internal/common"
)

var templates *template.Template

func renderTemplate(ctx *common.Context, tmpl string, pageData interface{}) {
	data := map[string]interface{}{
		"Title":       "osu!DiffCalc - web version",
		"Description": "An attempt to accurately compute beatmap difficulty and player ranking.",
		"LoadTime":    fmt.Sprintf("%.4f", time.Since(ctx.Start).Seconds()),
	}
	if pageData != nil {
		for k, v := range pageData.(map[string]interface{}) {
			data[k] = v
		}
	}

	err := templates.ExecuteTemplate(ctx.Response, tmpl, data)
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
