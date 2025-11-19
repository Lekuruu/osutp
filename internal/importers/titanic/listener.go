package titanic

import (
	"net/http"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/gorilla/websocket"
)

func (importer *TitanicImporter) ListenForServerUpdates(state *common.State) error {
	c, _, err := websocket.DefaultDialer.Dial(state.Config.Server.ApiEventsUrl, http.Header{
		"Authorization": []string{state.Config.Server.ApiEventsAuth},
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

		state.Logger.Logf("Received server update: %s", message)
	}
}
