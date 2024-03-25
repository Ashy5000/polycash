// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

var securityLevel = 0
var initialBlockDifficulty = uint64(100000)
var minimumBlockDifficulty = uint64(100000)
var maximumUint64 = ^uint64(0)
var blocksBeforeSpendable = 3

// Mining power is measured in difficulty points per minute (DPM).
const dpm = 1
const kdpm = 1000
const mdpm = 1000000
const gdpm = 1000000000
