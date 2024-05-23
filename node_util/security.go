// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

type SecurityLevel struct {
	Level                  int
	MinimumDifficulty      uint64
	InitialBlockDifficulty uint64
	BlocksBeforeReward     int
}

var securityLevels = []SecurityLevel{
	{0, 40000, 50000, 3},
	{1, 100000, 120000, 5},
	{2, 500000, 500000, 7},
}

func ApplySecurityLevel(level int) {
	for _, securityLevel := range securityLevels {
		if securityLevel.Level == level {
			InitialBlockDifficulty = securityLevel.InitialBlockDifficulty
			MinimumBlockDifficulty = securityLevel.MinimumDifficulty
			BlocksBeforeReward = securityLevel.BlocksBeforeReward
		}
	}
}
