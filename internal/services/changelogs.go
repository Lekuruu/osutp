package services

import (
	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
)

type ChangelogBatch struct {
	Date       string
	Changelogs []*database.Changelog
}

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

func FetchChangelogs(limit int, state *common.State) ([]*ChangelogBatch, error) {
	var logs []database.Changelog
	result := state.Database.Order("created_at DESC").Limit(limit).Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	var batches []*ChangelogBatch
	for i := range logs {
		entry := &logs[i]
		day := entry.Date()

		if len(batches) == 0 {
			batches = append(batches, &ChangelogBatch{
				Date:       day,
				Changelogs: []*database.Changelog{entry},
			})
			continue
		}

		previousBatch := batches[len(batches)-1]
		if previousBatch.Date != day {
			batches = append(batches, &ChangelogBatch{
				Date:       day,
				Changelogs: []*database.Changelog{entry},
			})
			continue
		}

		previousBatch.Changelogs = append(previousBatch.Changelogs, entry)
	}

	return batches, nil
}
