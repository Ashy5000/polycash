package main

import (
	"cryptocurrency/node_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransitionState(t *testing.T) {
	t.Run("It correctly transitions the state", func(t *testing.T) {
		// Arrange
		state := node_util.State{
			LegacyData: map[string][]byte{
				"123": []byte("321"),
			},
			LegacyContracts: map[uint64]node_util.Contract{
				0: {
					Contents: "",
					Parties:  nil,
					GasUsed:  0,
					Location: 0,
					Loaded:   false,
				},
			},
		}
		transition := node_util.StateTransition{
			LegacyUpdatedData: map[string][]byte{
				"321": []byte("123"),
			},
			LegacyNewContracts: map[uint64]node_util.Contract{
				1: {
					Contents: "",
					Parties:  nil,
					GasUsed:  0,
					Location: 0,
					Loaded:   false,
				},
			},
		}
		// Act
		newState := node_util.TransitionState(state, transition)
		// Assert
		for key, value := range transition.LegacyUpdatedData {
			state.LegacyData[key] = value
		}
		for key, value := range transition.LegacyNewContracts {
			state.LegacyContracts[key] = value
		}
		assert.Equal(t, state, newState)
	})
}
