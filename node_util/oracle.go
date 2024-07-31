// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type OracleType uint64

const ORACLE_RESPONSE_PREFIX = "ORACLE"

const (
	NilType OracleType = iota
)

type OracleRequest struct {
	Body []byte
	Type OracleType
}

type OracleResponse struct {
	Body    []byte
	Type    OracleType
	Request OracleRequest
}

type Oracle struct {
	Requests  []OracleRequest
	Responses []OracleResponse
}

func CalculateOracleResponse(request OracleRequest) OracleResponse {
	switch request.Type {
	case NilType:
		return OracleResponse{
			Body:    []byte("NULL"),
			Type:    0,
			Request: OracleRequest{},
		}
	}
	panic("Unknown oracle type.")
}

func (o *Oracle) Step() {
	if len(o.Requests) == 0 {
		return
	}
	for _, req := range o.Requests {
		o.Responses = append(o.Responses, CalculateOracleResponse(req))
	}
	o.Requests = []OracleRequest{}
}

func (o *Oracle) WriteResponses(transition StateTransition) {
	for _, res := range o.Responses {
		requestStr := []byte(fmt.Sprintf("%v", res.Request))
		requestHash := sha256.Sum256(requestStr)
		key := fmt.Sprintf("%s%s", ORACLE_RESPONSE_PREFIX, string(requestHash[:]))
		val, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		transition.UpdatedData[key] = val
	}
}
