package dispatcher

import "github.com/abhinavdahiya/go-messenger-bot"

// This interface defines any state for the bot
type State interface {
	// return the name of the state
	Name() string
	// force trasits the state
	Transit(string)
	// returns the next state
	Next() string
	// Actions
	Actions() (Action, Action)
	// Flush state to default
	Flush()
}

// This is the function that performs action
// on entering or leaving a particular state
type Action func(state State, c mbotapi.Callback, bot *mbotapi.BotAPI) error
