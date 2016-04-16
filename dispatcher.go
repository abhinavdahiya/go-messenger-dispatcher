package dispatcher

import (
	"errors"
	"log"

	"github.com/abhinavdahiya/go-messenger-bot"
)

type Storage interface {
	StoreState(mbotapi.User, State) error
	FetchState(mbotapi.User) (State, error)
}

type Dispatcher struct {
	States map[string]State
	Storage
	InitState string
	Debug     bool
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		States: make(map[string]State),
		Storage: &InMemory{
			Store: make(map[mbotapi.User]State),
		},
		InitState: "start",
		Debug:     false,
	}
}

func (d *Dispatcher) AddState(s State) {
	d.States[s.Name()] = s
}

func (d *Dispatcher) LoadState(name string) (State, error) {
	if s, ok := d.States[name]; ok {
		return s, nil
	}
	return nil, errors.New("Unknown State.")
}

func (d *Dispatcher) Process(c mbotapi.Callback, bot *mbotapi.BotAPI) error {
	// fetch the current state of the user
	curr, err := d.FetchState(c.Sender)
	if err != nil {
		// If no state found, initialize user to init state
		tmp := d.States[d.InitState]
		d.StoreState(c.Sender, tmp)
		curr = tmp
	}

	if d.Debug {
		log.Printf("[DEBUG] %#v", curr)
		log.Printf("[DEBUG] [CTX] %#v", Get(&curr))
	}

	var cLeave Action
	_, cLeave = curr.Actions()

	// exec Leave action for the state
	if cLeave != nil {
		err := cLeave(c, bot)
		if err != nil {
			return err
		}
	}

	// load next state
	curr.Transitor(c)
	ns := curr.Next()

	var next State
	next, err = d.LoadState(ns)
	var nEnter Action
	nEnter, _ = next.Actions()

	// load the next state
	ctx := Get(&curr)
	Set(&next, ctx)
	next.Flush()

	if d.Debug {
		log.Printf("[DEBUG] %#v", next)
		log.Printf("[DEBUG] [CTX] %#v", Get(&next))
	}

	if nEnter != nil {
		err = nEnter(c, bot)
		if err != nil {
			return err
		}
	}

	d.StoreState(c.Sender, next)

	if d.Debug {
		log.Printf("[DEBUG] %#v", d.Storage)
	}
	return nil
}
