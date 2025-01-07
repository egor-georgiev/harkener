package internal

import "context"

type State struct {
	Errors chan error
	Ctx    context.Context
}

func NewStateWithCancel() (*State, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	state := &State{
		Errors: make(chan error),
		Ctx:    ctx,
	}
	return state, cancel
}
