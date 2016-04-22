package dispatcher

import (
	"errors"
	"log"

	"github.com/abhinavdahiya/go-messenger-bot"
)

var (
	ErrNoCurrentState = errors.New("No current state for user")
	ErrUnknownState   = errors.New("Unknown state, state not known to dispatcher")
)

// Dispatcher stores all valid states
// and also the state of any particular user
// Also remember to set initstate , default to state
type Dispatcher struct {
	States    map[string]State
	Storage   map[mbotapi.User]State
	InitState string
	Debug     bool
}

// Creates a new Dispatcher with default settings
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		States:    make(map[string]State),
		Storage:   make(map[mbotapi.User]State),
		InitState: "start",
		Debug:     false,
	}
}

// Add state to the dispatcher
func (d *Dispatcher) AddState(states ...State) {
	for _, s := range states {
		d.States[s.Name] = s
	}
}

// Gets a state registered with dispatcher
func (d *Dispatcher) getState(name string) (State, error) {
	if s, ok := d.States[name]; ok {
		return CloneState(s), nil
	}
	return State{}, ErrUnknownState
}

// Returns current state of user
func (d *Dispatcher) GetUserState(u mbotapi.User) (State, error) {
	if s, ok := d.Storage[u]; ok {
		return s, nil
	}
	return State{}, ErrNoCurrentState
}

// Stores the current state of user
func (d *Dispatcher) StoreUserState(u mbotapi.User, s State) {
	if ref, ok := d.States[s.Name]; ok {
		clone := CloneState(ref)
		d.Storage[u] = clone
	}
}

func (d *Dispatcher) Process(c mbotapi.Callback, bot *mbotapi.BotAPI) error {
	if d.Debug {
		log.Printf("[DEBUG] (User: %v)Processing Starting =================================", c.Sender)
	}
	// fetch the current state of the user
	// if ErrNoCurrentState load init state to user
	curr, err := d.GetUserState(c.Sender)
	if err == ErrNoCurrentState {
		is, serr := d.getState(d.InitState)
		if serr != nil {
			return serr
		}
		d.StoreUserState(c.Sender, is)
		curr = is
	} else if err != nil {
		return err
	}

	if d.Debug {
		log.Printf("[DEBUG] (USER:%v) (STATE:%v)", c.Sender, curr)
		log.Printf("[DEBUG] [CTX] %#v", curr.GetCTX())
	}

	// Load leave action for curr state and
	// exec it
	var cl Action
	_, cl = curr.Actions()
	if cl != nil {
		err := cl(curr, c, bot)
		if err != nil {
			log.Printf("[ERROR] (Message: %v) (Error: %s)", c, err)
			return err
		}
	}

	// Find next state
	// If the next state is empty
	// move to initial state
	// Should you process the message through InitState??
	ns := curr.Next()
	if ns == "" {
		ts, _ := d.getState(d.InitState)
		var tl Action
		_, tl = ts.Actions()
		if tl != nil {
			tl(ts, c, bot)
		}
		ns = ts.Next()
	}

	if d.Debug {
		log.Printf("[DEBUG] (User: %v) (Next State: %s)", c.Sender, ns)
	}

	// Load next state
	// Move CTX
	// Run Enter action of next
	next, nerr := d.getState(ns)
	if nerr == nil {
		return err
	}

	if ctx := curr.GetCTX(); ctx != nil {
		next.SetCTX(ctx)
	}

	if d.Debug {
		log.Printf("[DEBUG] (USER:%v) (STATE:%v)", c.Sender, next)
		log.Printf("[DEBUG] [CTX] %#v", next.GetCTX())
	}

	var ne Action
	ne, _ = next.Actions()
	if ne != nil {
		err = ne(next, c, bot)
		if err != nil {
			log.Printf("[ERROR] (Message: %v) (Error: %s)", c, err)
			return err
		}
	}

	// Update new state of the user
	d.StoreUserState(c.Sender, next)

	return nil
}
