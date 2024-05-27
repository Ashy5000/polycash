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
	state := State{}
	for _, block := range Blockchain {
		state = TransitionState(state, block.Transition)
	}
	return state
}
