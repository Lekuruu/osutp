package titanic

import (
	"encoding/json"
	"net/http"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/gorilla/websocket"
)

type TitanicEvent struct {
	UserID int                    `json:"user_id"`
	Mode   int                    `json:"mode"`
	Type   UserActivityType       `json:"type"`
	Data   map[string]interface{} `json:"data"`
}

func (importer *TitanicImporter) ListenForServerUpdates(state *common.State) error {
	c, _, err := websocket.DefaultDialer.Dial(state.Config.Server.ApiEventsUrl, http.Header{
		"Authorization": []string{state.Config.Server.ApiAuth},
	})
	if err != nil {
		state.Logger.Logf("Error connecting to websocket: %v", err)
		return err
	}
	defer c.Close()

	state.Logger.Logf(
		"Listening for server updates on '%s'",
		state.Config.Server.ApiEventsUrl,
	)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			state.Logger.Log("Error reading websocket message:", err)
			continue
		}

		var event TitanicEvent
		if err := json.Unmarshal(message, &event); err != nil {
			state.Logger.Logf("Error unmarshaling server update: %v", err)
			continue
		}

		importer.handleServerEvent(state, event)
	}
}

func (importer *TitanicImporter) handleServerEvent(state *common.State, event TitanicEvent) {
	switch event.Type {
	case ActivityUserRegistration:
	case ActivityBeatmapLeaderboardRank:
		if event.Mode != 0 {
			return
		}
		state.Logger.Logf("Received server update for user %d (%d)", event.UserID, event.Type)
		go importer.ImportUser(event.UserID, state)
	case ActivityBeatmapStatusUpdated:
	case ActivityBeatmapUploaded:
	case ActivityBeatmapUpdated:
	case ActivityBeatmapRevived:
		beatmapsetID := event.Data["beatmapset_id"].(int)
		state.Logger.Logf("Received server update for beatmapset %d (%d)", beatmapsetID, event.Type)
		go importer.ImportBeatmapset(beatmapsetID, false, state)
	}
}
