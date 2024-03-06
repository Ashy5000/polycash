// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/dsa"
	"math/big"
	"strings"
)

func DecodePublicKey(keyString string) dsa.PublicKey {
	fields := strings.Split(keyString, "&")
	var y big.Int
	y.SetString(fields[0], 10)
	var p big.Int
	p.SetString(fields[1], 10)
	var q big.Int
	q.SetString(fields[2], 10)
	var g big.Int
	g.SetString(fields[3], 10)
	publicKey := dsa.PublicKey{
		Parameters: dsa.Parameters{
			P: &p,
			Q: &q,
			G: &g,
		},
		Y: &y,
	}
	return publicKey
}

func EncodePublicKey(key dsa.PublicKey) string {
	return key.Y.String() + "&" + key.P.String() + "&" + key.Q.String() + "&" + key.G.String()
}
