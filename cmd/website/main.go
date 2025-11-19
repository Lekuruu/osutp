package main

import (
	"log"
	"net/http"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/importers"
	"github.com/Lekuruu/osutp/internal/routes"
	"github.com/Lekuruu/osutp/internal/services"
)

func InitializeRoutes(server *common.Server) {
	server.Router.HandleFunc("/", server.ContextMiddleware(routes.Index)).Methods("GET")
	server.Router.HandleFunc("/info", server.ContextMiddleware(routes.Info)).Methods("GET")
	server.Router.HandleFunc("/scores", server.ContextMiddleware(routes.Scores)).Methods("GET")
	server.Router.HandleFunc("/players", server.ContextMiddleware(routes.Players)).Methods("GET")
	server.Router.HandleFunc("/banners", server.ContextMiddleware(routes.Banners)).Methods("GET")
	server.Router.HandleFunc("/beatmaps", server.ContextMiddleware(routes.Beatmaps)).Methods("GET")
	server.Router.HandleFunc("/changelog", server.ContextMiddleware(routes.Changelog)).Methods("GET")
	server.Router.HandleFunc("/banners/{pid:[0-9]+}", server.ContextMiddleware(routes.Banners)).Methods("GET")
	server.Router.HandleFunc("/banner/{pid:[0-9]+}", server.ContextMiddleware(routes.BannerImage)).Methods("GET")
	server.Router.HandleFunc("/players/{country:[a-zA-Z]{2}}", server.ContextMiddleware(routes.PlayersByCountry)).Methods("GET")
	// TODO: Implement dynamic strain graphs

	// Initialize static routes
	server.Router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/static/js/"))))
	server.Router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("web/static/css/"))))
	server.Router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("web/static/images/"))))
	server.Router.PathPrefix("/straingraph/").Handler(http.StripPrefix("/straingraph/", http.FileServer(http.Dir("web/static/straingraph/"))))
	server.Router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/robots.txt")
	})
	server.Router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/favicon.ico")
	})
}

func ResetIsUpdatingStatus(state *common.State) {
	players, err := services.FetchAllPlayers(state)
	if err != nil {
		log.Fatalf("Failed to fetch players: %v", err)
		return
	}

	for _, player := range players {
		if player.IsUpdating {
			log.Printf("Resetting updating status for player %d", player.ID)
			services.SetPlayerUpdatingStatus(player.ID, false, state)
		}
	}
}

func main() {
	state := common.NewState()
	if state == nil {
		return
	}

	log.SetFlags(0)
	log.SetOutput(state.Logger)

	importer, err := importers.NewImporter(state.Config)
	if err != nil {
		log.Fatalf("Failed to create importer: %v", err)
		return
	}
	state.Extensions["importer"] = importer

	// Server might have restarted while players were being updated
	// without having their status reset
	go ResetIsUpdatingStatus(state)

	server := common.NewServer(
		state.Config.Web.Host,
		state.Config.Web.Port,
		"osu!tp",
		state,
	)
	InitializeRoutes(server)
	server.Serve()
}
