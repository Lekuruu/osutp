package titanic

import (
	"encoding/json"
	"fmt"
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
	attempt := 0

	for {
		connected, connectedFor, err := importer.listenForUpdatesUntilDisconnect(state)
		if connected && connectedFor >= time.Minute {
			attempt = 0
		}
		if err != nil {
			state.Logger.Logf("Websocket event listener disconnected: %v", err)
		}

		delay := websocketReconnectDelay(attempt)
		state.Logger.Logf("Reconnecting websocket event listener in %s", delay)
		time.Sleep(delay)
		attempt++
	}
}

func (importer *TitanicImporter) listenForUpdatesUntilDisconnect(state *common.State) (connected bool, connectedFor time.Duration, err error) {
	var connectedAt time.Time
	defer func() {
		if r := recover(); r != nil {
			if connected {
				connectedFor = time.Since(connectedAt)
			}
			err = fmt.Errorf("panic while listening for server updates: %v", r)
			state.Logger.Logf("Recovered from websocket listener panic: %v", r)
		}
	}()
	state.Logger.Logf("Connecting websocket event listener to '%s'", state.Config.Server.ApiEventsUrl)

	c, _, err := websocket.DefaultDialer.Dial(state.Config.Server.ApiEventsUrl, http.Header{
		"Authorization": []string{state.Config.Server.ApiAuth},
	})
	if err != nil {
		return false, 0, fmt.Errorf("failed to connect websocket: %w", err)
	}

	connected = true
	connectedAt = time.Now()
	defer func() {
		if closeErr := c.Close(); closeErr != nil {
			state.Logger.Logf("Error closing websocket event listener: %v", closeErr)
		}
		state.Logger.Log("Websocket event listener connection closed")
	}()

	state.Logger.Logf(
		"Websocket event listener connected to '%s'",
		state.Config.Server.ApiEventsUrl,
	)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			connectedFor = time.Since(connectedAt)
			return connected, connectedFor, fmt.Errorf("failed to read websocket message after %s: %w", connectedFor.Round(time.Second), err)
		}

		var event TitanicEvent
		if err := json.Unmarshal(message, &event); err != nil {
			state.Logger.Logf("Error unmarshaling server update: %v", err)
			continue
		}

		go importer.handleServerEvent(state, event)
	}
}

func websocketReconnectDelay(attempt int) time.Duration {
	if attempt > 5 {
		attempt = 5
	}
	return backoffDelay(attempt)
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
