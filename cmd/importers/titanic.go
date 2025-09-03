package main

import (
	"fmt"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/importers/titanic"
)

func main() {
	state := common.NewState()
	if state == nil {
		return
	}

	err := titanic.ImportBeatmapsByDifficulty(0, state)
	if err != nil {
		fmt.Printf("Error occurred while importing beatmaps: %v\n", err)
	}
}
