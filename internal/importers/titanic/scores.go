package titanic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
	"github.com/Lekuruu/osutp-web/internal/services"
	"gorm.io/gorm"
)

type ScoreCollectionModel struct {
	Total  int          `json:"total"`
	Scores []ScoreModel `json:"scores"`
}

type ScoreModel struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	BeatmapID     int       `json:"beatmap_id"`
	SubmittedAt   string    `json:"submitted_at"`
	Mode          int       `json:"mode"`
	StatusPP      int       `json:"status_pp"`
	StatusScore   int       `json:"status_score"`
	ClientVersion int       `json:"client_version"`
	PP            float64   `json:"pp"`
	PPv1          float64   `json:"ppv1"`
	Acc           float64   `json:"acc"`
	TotalScore    int       `json:"total_score"`
	MaxCombo      int       `json:"max_combo"`
	Mods          int       `json:"mods"`
	Perfect       bool      `json:"perfect"`
	Passed        bool      `json:"passed"`
	Pinned        bool      `json:"pinned"`
	Count300      int       `json:"n300"`
	Count100      int       `json:"n100"`
	Count50       int       `json:"n50"`
	CountMiss     int       `json:"nMiss"`
	CountGeki     int       `json:"nGeki"`
	CountKatu     int       `json:"nKatu"`
	Grade         string    `json:"grade"`
	ReplayViews   int       `json:"replay_views"`
	Failtime      *int      `json:"failtime,omitempty"`
	User          UserModel `json:"user"`
}

func (score *ScoreModel) ToSchema() *database.Score {
	createdAt, err := time.Parse("2006-01-02T15:04:05", score.SubmittedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}

	return &database.Score{
		ID:         score.ID,
		BeatmapID:  score.BeatmapID,
		PlayerID:   score.UserID,
		TotalScore: score.TotalScore,
		MaxCombo:   score.MaxCombo,
		Mods:       uint32(score.Mods),
		FullCombo:  score.Perfect,
		Grade:      score.Grade,
		Accuracy:   score.Acc,
		Amount300:  score.Count300,
		Amount100:  score.Count100,
		Amount50:   score.Count50,
		AmountGeki: score.CountGeki,
		AmountKatu: score.CountKatu,
		AmountMiss: score.CountMiss,
		CreatedAt:  createdAt,
	}
}

func ImportOrUpdateLeaderboards(beatmaps []*database.Beatmap, state *common.State) {
	for _, beatmap := range beatmaps {
		scores, err := PerformLeaderboardRequest(beatmap.ID, 0, state)
		if err != nil {
			fmt.Println("Failed to fetch leaderboard for beatmap", beatmap.ID, ":", err)
			continue
		}

		if err := ProcessScores(scores, beatmap, state); err != nil {
			fmt.Println("Failed to process scores for beatmap", beatmap.ID, ":", err)
			continue
		}
	}
}

func ProcessScores(scores []ScoreModel, beatmap *database.Beatmap, state *common.State) error {
	for _, score := range scores {
		if score.Mode != 0 {
			// Only process osu! standard scores
			continue
		}

		scoreExists, err := services.ScoreExists(score.ID, state)
		if err != nil {
			fmt.Println("Failed to check if score exists:", err)
			continue
		}
		if scoreExists {
			continue
		}

		schema := score.ToSchema()
		difficulty, err := beatmap.DifficultyCalculationResult(schema.DifficultyMods())
		if err != nil {
			fmt.Println("Failed to get difficulty calculation result:", err)
			continue
		}

		request := schema.CalculationRequest(difficulty)
		response, err := request.Perform(state.Config.TpServiceUrl)
		if err != nil {
			fmt.Println("Failed to calculate performance:", err)
			continue
		}

		schema.TotalTp = response.Total
		schema.AimTp = response.Aim
		schema.SpeedTp = response.Speed
		schema.AccTp = response.Acc

		user, err := services.FetchUserById(score.UserID, state)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println("Failed to fetch user:", err)
			continue
		}

		if user == nil {
			user = score.User.ToSchema()
			if err := services.CreateUser(user, state); err != nil {
				fmt.Println("Failed to create user:", err)
				continue
			}
		}

		personalBest, err := services.FetchPersonalBestScore(user.ID, beatmap.ID, state)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println("Failed to fetch personal best score:", err)
			continue
		}

		if personalBest != nil && personalBest.TotalTp > schema.TotalTp {
			continue
		}

		// Delete old personal best
		if personalBest != nil {
			if err := services.DeleteScore(personalBest.ID, state); err != nil {
				fmt.Println("Failed to delete old personal best score:", err)
				continue
			}
		}

		if err := services.CreateScore(schema, state); err != nil {
			fmt.Println("Failed to create score:", err)
			continue
		}

		fmt.Printf("Imported score from '%s' on beatmap '%s' with %.2ftp\n", user.Name, beatmap.FullName(), schema.TotalTp)
	}
	return nil
}

func PerformLeaderboardRequest(beatmapId int, offset int, state *common.State) ([]ScoreModel, error) {
	url := fmt.Sprintf("%s/beatmaps/%d/scores?offset=%d", state.Config.Server.ApiUrl, beatmapId, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch leaderboard: %s", resp.Status)
	}

	var scores ScoreCollectionModel
	if err := json.NewDecoder(resp.Body).Decode(&scores); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return scores.Scores, nil
}
