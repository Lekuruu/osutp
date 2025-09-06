package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/services"
	"github.com/xeonx/timeago"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var templates *template.Template
var printer = message.NewPrinter(language.English)

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
		"Server":      ctx.State.Config.Server,
		"Query":       ctx.Request.URL.Query(),
		"Printer":     printer,
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
	funcs := template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"sub":   func(a, b int) int { return a - b },
		"mod":   func(a, b int) int { return a % b },
		"mul":   func(a, b int) int { return a * b },
		"div":   func(a, b int) int { return a / b },
		"lower": func(s string) string { return strings.ToLower(s) },
		"upper": func(s string) string { return strings.ToUpper(s) },
		"query": func(name, defaultValue string, q url.Values) string {
			value := q.Get(name)
			if value == "" {
				return defaultValue
			}
			return value
		},
	}

	var err error
	templates, err = template.New("").Funcs(funcs).ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}
