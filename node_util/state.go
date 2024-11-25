// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

type State struct {
	// Legacy properties are stored as maps for backwards compatibility
	LegacyData      map[string][]byte
	LegacyContracts map[uint64]Contract
	// Post-zen properties are stored as merkle trees
	ZenData      []MerkleNode
	ZenContracts []MerkleNode
}

type StateTransition struct {
	LegacyUpdatedData  map[string][]byte
	LegacyNewContracts map[uint64]Contract
	ZenUpdatedData     []MerkleNode
	ZenNewContracts    []MerkleNode
}

func TransitionState(state State, transition StateTransition) State {
	for key, value := range transition.LegacyUpdatedData {
		if len(value) == 0 || value == nil {
			continue
		}
		state.LegacyData[key] = value
	}
	for key, value := range transition.LegacyNewContracts {
		state.LegacyContracts[key] = value
	}
	state.ZenData = Merge(state.ZenData, transition.ZenUpdatedData)
	state.ZenContracts = Merge(state.ZenContracts, transition.ZenNewContracts)
	return state
}

func CalculateCurrentState() State {
	state := State{
		LegacyData:      make(map[string][]byte),
		LegacyContracts: make(map[uint64]Contract),
		ZenData:         make([]MerkleNode, 0),
		ZenContracts:    make([]MerkleNode, 0),
	}
	for _, block := range Blockchain {
		state = TransitionState(state, block.Transition)
	}
	return state
}

func GetFromState(location string) []byte {
	state := CalculateCurrentState()
	// Try to get from merkle
	val, ok := GetValue(state.ZenData, location)
	if ok {
		return val
	}
	// Fallback to legacy
	val, ok = state.LegacyData[location]
	if ok {
		return val
	}
	return []byte{}
}

// GetPendingState is for legacy versions only. Post-zen VMs process the entire batch in one ZK proof
func GetPendingState() map[string][]byte {
	res := make(map[string][]byte)
	for _, subTransition := range NextTransitions {
		for location, data := range subTransition.LegacyUpdatedData {
			res[location] = data
		}
	}
	return res
}
