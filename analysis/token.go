package analysis

import (
	. "cryptocurrency/node_util"
)

func GetNumTokensMinted() int64 {
	var result float64
	for i, block := range Blockchain {
		if i > 0 {
			lastBlock := Blockchain[i-1]
			result += float64(len(block.TimeVerifiers)-len(lastBlock.TimeVerifiers)) * 0.1
		}
		minerCount := GetMinerCount(i)
		reward := CalculateBlockReward(minerCount, i)
		result += reward
	}
	if len(Blockchain) > 50 {
		result -= float64(int(GetMinerCount(len(Blockchain))) * BlocksBeforeReward) // First n blocks for each miner don't have a reward
	}
	return int64(result)
}
