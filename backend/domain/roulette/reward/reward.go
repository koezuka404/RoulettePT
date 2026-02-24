package reward

type Reward struct {
	Points int
	Rate   int
}

var Rewards = []Reward{
	{10, 50},
	{50, 30},
	{100, 15},
	{500, 5},
}
