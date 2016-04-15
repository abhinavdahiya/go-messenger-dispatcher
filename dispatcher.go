package dispatcher

import (
	"errors"

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
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		States: make(map[string]State),
		Storage: &InMemory{
			Store: make(map[mbotapi.User]State),
		},
		InitState: "start",
	}
}

func (d *Dispatcher) AddState(s State) {
	d.States[s.Name] = s
}

func (d *Dispatcher) LoadState(name string) (State, error) {
	if s, ok := d.States[name]; ok {
		return s, nil
	}
	return State{}, errors.New("Unknown State.")
}

func (d *Dispatcher) Process(c mbotapi.Callback, bot *mbotapi.BotAPI) error {
	// fetch the current state of the user
	curr, err := d.FetchState(c.Sender)
	if err != nil {
		// If no state found, initialize user to init state
		tmp := d.States[d.InitState]
		d.StoreState(c.Sender, tmp)
		return err
	}

	// load next state
	var ns string
	if curr.IsMoved {
		ns = curr.Chain
	}
	ns, err = curr.Test(c)
	if err != nil {
		return err
	}

	var next State
	next, err = d.LoadState(ns)

	// exec Leave action for the state
	if curr.Leave != nil {
		err := curr.Leave(c, bot)
		if err != nil {
			return err
		}
	}

	// load the next state
	ctx := curr.Data()
	next.SetData(ctx)
	err = next.Enter(c, bot)
	if err != nil {
		return err
	}

	d.StoreState(c.Sender, next)
	return nil
}
