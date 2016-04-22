package dispatcher

import "github.com/abhinavdahiya/go-messenger-bot"

// This struct defines any state for the bot
type State struct {
	Name         string
	Enter, Leave Action
	IsMoved      bool
	Chain        string
}

// Creates a new state with `name`
func NewState(name string) State {
	s := State{
		Name:    name,
		IsMoved: false,
		Chain:   "",
	}

	return s
}

// Returns a new instance of the same state
func CloneState(orig State) State {
	clone := State{
		Name:    orig.Name,
		IsMoved: false,
		Chain:   "",
		Enter:   orig.Enter,
		Leave:   orig.Leave,
	}
	return clone
}

// Transits the the current state to new state
func (s *State) Transit(new string) {
	s.IsMoved = true
	s.Chain = new
}

// Gives the Next state
func (s *State) Next() string {
	if s.IsMoved {
		return s.Chain
	}
	return ""
}

// Add state handlers for state
func (s *State) SetActions(e, l Action) {
	s.Enter = e
	s.Leave = l
}

// Returns enter and leave actions of a state
func (s *State) Actions() (Action, Action) {
	return s.Enter, s.Leave
}

// This is the function that performs action
// on entering or leaving a particular state
type Action func(state *State, c mbotapi.Callback, bot *mbotapi.BotAPI) error
