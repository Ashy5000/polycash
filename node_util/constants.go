// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

// Difficulty
var InitialBlockDifficulty = uint64(50000)
var MinimumBlockDifficulty = uint64(50000)
var MaximumUint64 = ^uint64(0)

// Finality
const BlocksUntilFinality = 3

// Rewards
var BlocksBeforeReward = 3
var BlockReward = 1.0
var TransactionFee = 0.0001
var BodyFeePerByte = 0.000001
var GasPrice = 0.000001

// Mining power is measured in difficulty points per minute (DPM).
const Dpm = 1
const Kdpm = 1000
const Mdpm = 1000000
const Gdpm = 1000000000
