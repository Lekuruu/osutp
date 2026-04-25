package titanic

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
	"github.com/Lekuruu/osutp/internal/updaters"
	"github.com/gorilla/websocket"
)

type TitanicEvent struct {
	UserID int              `json:"user_id"`
	Mode   int              `json:"mode"`
	Type   UserActivityType `json:"type"`
	Data   map[string]any   `json:"data"`
}

func (importer *TitanicImporter) ListenForServerUpdates(state *common.State) error {
	defer func() {
		if r := recover(); r != nil {
			state.Logger.Logf("Recovered from panic: %v", r)
			go importer.ListenForServerUpdates(state)
		}
	}()

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

		go importer.handleServerEvent(state, event)
	}
}

func (importer *TitanicImporter) handleServerEvent(state *common.State, event TitanicEvent) {
	switch event.Type {
	case ActivityUserRegistration:
		state.Logger.Logf("Received server update for user %d (%d)", event.UserID, event.Type)

		if _, err := importer.ImportUser(event.UserID, state); err != nil {
			state.Logger.Logf("Error importing user %d from event: %v", event.UserID, err)
			return
		}
		importer.onServerEventFinished(state)
	case ActivityBeatmapLeaderboardRank:
		if event.Mode != 0 {
			return
		}
		state.Logger.Logf("Received server update for user %d (%d)", event.UserID, event.Type)

		if _, err := importer.ImportUser(event.UserID, state); err != nil {
			state.Logger.Logf("Error importing user %d from event: %v", event.UserID, err)
		}

		if beatmapIdFloat, ok := event.Data["beatmap_id"].(float64); !ok {
			state.Logger.Logf("Error: beatmap_id not found or invalid type in event data")
		} else {
			beatmapId := int(beatmapIdFloat)
			if _, err := importer.ImportBeatmap(beatmapId, true, state); err != nil {
				state.Logger.Logf("Error importing beatmap %d from event: %v", beatmapId, err)
			}
		}

		importer.onServerEventFinished(state)
	case ActivityBeatmapStatusUpdated,
		ActivityBeatmapUploaded,
		ActivityBeatmapUpdated,
		ActivityBeatmapRevived:
		beatmapsetIDFloat, ok := event.Data["beatmapset_id"].(float64)
		if !ok {
			state.Logger.Logf("Error: beatmapset_id not found or invalid type in event data")
			return
		}
		beatmapsetID := int(beatmapsetIDFloat)
		state.Logger.Logf("Received server update for beatmapset %d (%d)", beatmapsetID, event.Type)

		if _, err := importer.ImportBeatmapset(beatmapsetID, false, state); err != nil {
			state.Logger.Logf("Error importing beatmapset %d from event: %v", beatmapsetID, err)
			return
		}

		importer.onServerEventFinished(state)
	default:
		return
	}
}

func (importer *TitanicImporter) onServerEventFinished(state *common.State) {
	if err := updaters.UpdatePlayerRatings(state); err != nil {
		state.Logger.Logf("Error updating player ratings: %v", err)
		return
	}

	// This will be displayed in the header: "updated daily - last update: <...>"
	services.UpdatePageLastUpdated("players", time.Now(), state)
}
