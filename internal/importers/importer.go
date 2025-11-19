package importers

import (
	"fmt"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/importers/titanic"
)

type Importer interface {
	ImportUser(userID int, state *common.State) (*database.Player, error)
	ImportBeatmap(beatmapID int, importLeaderboard bool, state *common.State) (*database.Beatmap, error)
	ImportScore(scoreID int, state *common.State) (*database.Score, error)

	ImportUsersFromRankings(page int, state *common.State) (int, error)
	ImportBeatmapsByDifficulty(page int, state *common.State) (int, error)
	ImportBeatmapsByDate(page int, state *common.State) (int, error)

	EnqueueUserUpdate(userID int, state *common.State) error
	ListenForServerUpdates(state *common.State) error
}

func NewImporter(config *common.Config) (Importer, error) {
	switch config.Server.Type {
	case "titanic":
		return titanic.NewTitanicImporter(config.Server.WebUrl, config.Server.ApiUrl), nil
	default:
		return nil, fmt.Errorf("unknown importer type: %s", config.Server.Type)
	}
}
