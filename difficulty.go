package main

import (
	"math"
	"time"
)

func GetDifficulty(last_time time.Duration, last_difficulty uint64) uint64 {
	// The target time for a block is 1 minute.
	// The difficulty is adjusted on a per-miner, per-block basis.
	// To give faster miners a (small) advantage, the difficulty is divided by the result of a modified sigmoid function.
	// It looks like this:
	// 1 / (1 + e^(-(x-mdpm)/mdpm))
	// Where x is the previous time times the previous difficulty, and 1 mdpm is  1000000 (1 million) difficulty points per minute.
	difficultyBeforeAdjustment := last_difficulty * uint64(60) / uint64(last_time.Seconds())
	x := last_time.Minutes() * float64(last_difficulty)
	adjustment := (1 / (1 + math.Pow(math.E, -(1/mdpm)*(x-mdpm)))) + 0.5
	difficultyAfterAdjustment := float64(difficultyBeforeAdjustment) / adjustment
	return uint64(difficultyAfterAdjustment)
}
