// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

type OracleQueryType uint64

const (
	NilType OracleQueryType = iota
)

type OracleQuery struct {
	Body []byte
	Type OracleQueryType
}

type OracleResponse struct {
	Body []byte
}

func CalculateOracleResponse(query OracleQuery) OracleResponse {
	switch query.Type {
	case NilType:
		return OracleResponse{
			Body: []byte("NULL"),
		}
	}
	panic("Unknown oracle type.")
}
