package banners

import "fmt"

type Player interface {
	Username() string
	Country() string
	CountryCode() string
	CountryFlagLocation() string
	GlobalRank() int
	CountryRank() int
	CountryRankOrdinal() string
	Tp() int
}

type BasePlayer struct {
	username    string
	country     string
	countryCode string
	countryFlag string
	globalRank  int
	countryRank int
	tp          int
}

func (p *BasePlayer) Username() string {
	return p.username
}

func (p *BasePlayer) Country() string {
	return p.country
}

func (p *BasePlayer) CountryCode() string {
	return p.countryCode
}

func (p *BasePlayer) CountryFlagLocation() string {
	return p.countryFlag
}

func (p *BasePlayer) GlobalRank() int {
	return p.globalRank
}

func (p *BasePlayer) CountryRank() int {
	return p.countryRank
}

func (p *BasePlayer) Tp() int {
	return p.tp
}

func (p *BasePlayer) CountryRankOrdinal() string {
	rank := p.CountryRank()

	if rank >= 11 && rank <= 13 {
		return fmt.Sprintf("%dth", rank)
	}

	switch rank % 10 {
	case 1:
		return fmt.Sprintf("%dst", rank)
	case 2:
		return fmt.Sprintf("%dnd", rank)
	case 3:
		return fmt.Sprintf("%drd", rank)
	default:
		return fmt.Sprintf("%dth", rank)
	}
}

func NewPlayer(username, country, countryCode, countryFlag string, globalRank, countryRank, tp int) Player {
	return &BasePlayer{
		username:    username,
		country:     country,
		countryCode: countryCode,
		countryFlag: countryFlag,
		globalRank:  globalRank,
		countryRank: countryRank,
		tp:          tp,
	}
}
