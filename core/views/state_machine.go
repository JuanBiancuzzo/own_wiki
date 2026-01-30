package views

import (
	"fmt"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

type StateMachine struct {
	CurrentState  string
	ExecutedState string
	AllStates     []string
}

func NewStateMachine(inicialState string, states ...string) *StateMachine {
	return &StateMachine{
		CurrentState:  inicialState,
		ExecutedState: inicialState,
		AllStates:     append(states, inicialState),
	}
}

func (sm *StateMachine) Run(sCtx *s.SceneCtx, stateFuntions map[string]func(sCtx *s.SceneCtx) string) error {
	hasState := false
	for _, state := range sm.AllStates {
		if hasState = sm.CurrentState == state; hasState {
			break
		}
	}
	if !hasState {
		return fmt.Errorf("The new state %q is not in the states define at initialization", sm.CurrentState)
	}

	if function, ok := stateFuntions[sm.CurrentState]; !ok {
		return fmt.Errorf("The function for the current state %q, is not define", sm.CurrentState)

	} else {
		sm.ExecutedState = sm.CurrentState
		sm.CurrentState = function(sCtx)
		return nil
	}
}

func (sm *StateMachine) State() string {
	return sm.ExecutedState
}
