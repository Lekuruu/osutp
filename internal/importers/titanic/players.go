package titanic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
)

func (importer *TitanicImporter) ImportUser(userID int, state *common.State) (*database.Player, error) {
	user, err := importer.fetchUserById(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return importer.importUserFromModel(*user, state)
}

func (importer *TitanicImporter) ImportUsersFromRankings(page int, state *common.State) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (importer *TitanicImporter) importUserFromModel(user UserModel, state *common.State) (*database.Player, error) {
	// Check for existing player entry
	userEntry, err := services.FetchPlayerById(user.ID, state)
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	if userEntry != nil {
		// TODO: Update existing player data if necessary
		return userEntry, nil
	}

	userEntry = user.ToSchema()
	if err := services.PlayerUser(userEntry, state); err != nil {
		return nil, err
	}

	return userEntry, nil
}

func (importer *TitanicImporter) fetchUserById(userID int) (*UserModel, error) {
	url := fmt.Sprintf("%s/users/%d", importer.ApiUrl, userID)
	resp, err := http.Get(url)
	if err != nil {
		// Check for any rate limit errors and wait if needed
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			time.Sleep(time.Second * 60)
			return importer.fetchUserById(userID)
		}
		return nil, err
	}
	defer resp.Body.Close()

	var user UserModel
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &user, nil
}
