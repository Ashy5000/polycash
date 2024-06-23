// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

type State struct {
	Data map[string][]byte
}

type StateTransition struct {
	UpdatedData map[string][]byte
}

func TransitionState(state State, transition StateTransition) State {
	for key, value := range transition.UpdatedData {
		state.Data[key] = value
	}
	return state
}

func CalculateCurrentState() State {
	state := State{
		Data: make(map[string][]byte),
	}
	for _, block := range Blockchain {
		state = TransitionState(state, block.Transition)
	}
	return state
}
