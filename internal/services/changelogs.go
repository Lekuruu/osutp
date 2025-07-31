package services

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
)

func CreateChangelog(changelog *database.Changelog, state *common.State) error {
	result := state.Database.Create(changelog)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FetchChangelog(changelogId int, state *common.State) (*database.Changelog, error) {
	changelog := &database.Changelog{}
	result := state.Database.First(changelog, changelogId)
	if result.Error != nil {
		return nil, result.Error
	}
	return changelog, nil
}

func FetchChangelogs(limit int, state *common.State) (map[string][]*database.Changelog, error) {
	var logs []database.Changelog
	result := state.Database.Order("created_at DESC").Limit(limit).Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	reversed := make([]database.Changelog, len(logs))
	for i, entry := range logs {
		reversed[len(logs)-1-i] = entry
	}

	grouped := make(map[string][]*database.Changelog)
	for i := range reversed {
		entry := &reversed[i]
		day := entry.Date()
		grouped[day] = append(grouped[day], entry)
	}

	return grouped, nil
}
