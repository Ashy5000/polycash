package node_util

import "math"

func CalculateBlockReward(minerCount int64) float64 {
	// The more miners, the less reward
	// This is designed to prevent miners from forking their hash power to get more rewards
	r := math.Pow(10, 22)
	p := 0.95
	reward := math.Pow(p, float64(minerCount)) * r
	if reward > 10 {
		reward = 10
	}
	return reward
}
