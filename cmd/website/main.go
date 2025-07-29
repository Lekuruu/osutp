package main

import (
	"net/http"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/routes"
)

func InitializeRoutes(server *common.Server) {
	server.Router.HandleFunc("/", server.ContextMiddleware(routes.Index)).Methods("GET")
	server.Router.HandleFunc("/players", server.ContextMiddleware(routes.Players)).Methods("GET")

	// Initialize static routes
	server.Router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/static/js/"))))
	server.Router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("web/static/css/"))))
	server.Router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("web/static/images/"))))
	server.Router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/robots.txt")
	})
	server.Router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/favicon.ico")
	})
}

func main() {
	state := common.NewState()
	if state == nil {
		return
	}

	server := common.NewServer(
		state.Config.Web.Host,
		state.Config.Web.Port,
		"osu!tp",
		state,
	)
	InitializeRoutes(server)
	server.Serve()
}
