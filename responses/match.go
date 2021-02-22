package responses

// Match - Struct for a Match's json response
type Match struct {
	Day       uint
	Opponent  uint
	Concluded bool
	Won       bool
}
