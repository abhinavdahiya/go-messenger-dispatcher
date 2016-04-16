package dispatcher

import (
	"errors"

	"github.com/abhinavdahiya/go-messenger-bot"
)

type InMemory struct {
	Store map[mbotapi.User]State
}

func (i *InMemory) StoreState(u mbotapi.User, s State) error {
	i.Store[u] = s
	return nil
}

func (i *InMemory) FetchState(u mbotapi.User) (State, error) {
	if s, ok := i.Store[u]; ok {
		return s, nil
	}
	return nil, errors.New("State not found")
}
