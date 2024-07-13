// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"math"
	"time"
)

func GetDifficulty(lastTime time.Duration, lastDifficulty uint64) uint64 {
	// The target time for a block is 1 minute.
	// The difficulty is adjusted on a per-miner, per-block basis.
	// To give faster miners a (small) advantage, the difficulty is divided by the result of a modified sigmoid function.
	// It should be noted that miners with mining rates past a certain point will be disadvantaged.
	// It looks like this:
	// 1 / (1 + e^(-(x-mdpm)/mdpm))
	// Where x is the previous time times the previous difficulty, and 1 mdpm is  1000000 (1 million) difficulty points per minute.
	difficultyBeforeAdjustment := lastDifficulty * uint64(60) / uint64(lastTime.Seconds())
	x := lastTime.Minutes() * float64(lastDifficulty)
	var adjustment float64
	adjustment = (1 / (1 + math.Pow(math.E, -(1/(10*Kdpm))*(x-Mdpm)))) + 0.5
	difficultyAfterAdjustment := float64(difficultyBeforeAdjustment) / adjustment
	difficultyUint64 := uint64(difficultyAfterAdjustment)
	if difficultyUint64 > MinimumBlockDifficulty {
		return difficultyUint64
	}
	return MinimumBlockDifficulty
}
