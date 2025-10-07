# osutp

This project provides a recreation of the osu!tp algorithm and the now-defunct [osu!tp](https://web.archive.org/web/20131208212150/http://osutp.net/players) website. It makes use of [Tom94](https://github.com/Tom94)'s [AiModtpDifficultyCalculator](https://github.com/Tom94/AiModtpDifficultyCalculator), released back in 2013.

This project serves two purposes:

1. **As a Go package** - Calculate difficulty and performance using the osu!tp algorithm in your own projects
2. **As a web frontend** - Run a full recreation of the original osu!tp website

---

## Using as a Package

Import the packages in your Go project:

```go
import (
    "github.com/Lekuruu/osutp/pkg/tp"
    osu "github.com/natsukagami/go-osu-parser"
)
```

Calculate difficulty for a beatmap:

```go
beatmap, err := osu.ParseFile("beatmap.osu")
if err != nil {
    panic(err)
}

mods := tp.Hidden | tp.HardRock
difficulty := tp.CalculateDifficulty(&beatmap, mods)

fmt.Printf("Star Rating: %.2f*\n", difficulty.StarRating)
fmt.Printf("Aim: %.2f*, Speed: %.2f*\n", difficulty.AimStars, difficulty.SpeedStars)
```

Calculate performance for a score:

```go
score := &tp.Score{
    Amount300:  500,
    Amount100:  50,
    Amount50:   0,
    AmountMiss: 0,
    MaxCombo:   727,
    Mods:       tp.Hidden | tp.HardRock,
}

pp := tp.CalculatePerformance(difficulty, score)
fmt.Printf("Performance: %.2fpp\n", pp)

// You can also use `tp.NewScoreFromReplay` to load a score from a replay
// This will require the "github.com/robloxxa/go-osr" package
```

## Running the Website

The website provides a full recreation of the original osu!tp interface, including player rankings, beatmap listings, and banner generation.
The only requirement is to have [go](https://go.dev/) installed.

### Setup

Start off by cloning the repository onto your machine:

```bash
git clone https://github.com/Lekuruu/osutp.git
cd osutp
```

The server can be configured using environment variables. The default values are meant to be used with the Titanic! private server, but should work out of the box either way.

Finally, run the website:

```bash
go run cmd/website/main.go
```

The website will be available at `http://localhost:8080` by default.

### Database

The SQLite database is automatically created at `./.data/osutp.db` on the first run. You'll need to import data using the provided importers (see `cmd/importers/`) to populate the database with beatmaps, players, and scores.
