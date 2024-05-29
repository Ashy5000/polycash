package node_util

type State struct {
	Data map[uint64][]byte
}

type StateTransition struct {
	UpdatedData map[uint64][]byte
}

func TransitionState(state State, transition StateTransition) State {
	for key, value := range transition.UpdatedData {
		state.Data[key] = value
	}
	return state
}

func CalculateCurrentState() State {
	state := State{
		Data: make(map[uint64][]byte),
	}
	for _, block := range Blockchain {
		state = TransitionState(state, block.Transition)
	}
	return state
}
