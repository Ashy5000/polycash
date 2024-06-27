// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"fmt"
	"strconv"
	"strings"
)

func DecodePublicKey(keyString string) PublicKey {
	key := PublicKey{
		Y: []byte(""),
	}
	for _, ps := range strings.Split(strings.Trim(keyString, "[]"), " ") {
		pi, err := strconv.ParseUint(ps, 10, 8)
		if err != nil {
			panic(err)
		}
		key.Y = append(key.Y, byte(pi))
	}
	return key
}

func EncodePublicKey(key PublicKey) string {
	result := fmt.Sprintf("%v", key.Y)
	return result
}
