package responses

// Rikishi - Struct for a Rikishi's json response
type Rikishi struct {
	ID      uint
	Name    string
	Avatar  string
	Rank    string
	Wins    uint           `gorm:"-"`
	Losses  uint           `gorm:"-"`
	Matches map[uint]Match `gorm:"-"`
}
