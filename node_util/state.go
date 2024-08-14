// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

type State struct {
	Data      map[string][]byte
	Contracts map[uint64]Contract
}

type StateTransition struct {
	UpdatedData  map[string][]byte
	NewContracts map[uint64]Contract
}

func TransitionState(state State, transition StateTransition) State {
	for key, value := range transition.UpdatedData {
    if len(value) == 0 || value == nil {
      continue
    }
		state.Data[key] = value
	}
	for key, value := range transition.NewContracts {
		state.Contracts[key] = value
	}
	return state
}

func CalculateCurrentState() State {
	state := State{
		Data: make(map[string][]byte),
    Contracts: make(map[uint64]Contract),
	}
	for _, block := range Blockchain {
		state = TransitionState(state, block.Transition)
	}
	return state
}

func GetFromState(location string, state State) []byte {
  val, ok := state.Data[location];
  if ok {
    return val;
  }
  return []byte{};
}
