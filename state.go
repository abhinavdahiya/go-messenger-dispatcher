package dispatcher

import (
	"errors"
	"regexp"

	"github.com/abhinavdahiya/go-messenger-bot"
)

var (
	ErrNoTransitorFound = errors.New("No transitor found, Plz set fallback rule to prevent this.")
)

// This struct defines any state for the bot
// corresponding to any user
// This also stores transitors for all possible
// moves from this state
//
// make sure the state Name is unique
type State struct {
	Name          string
	Enter, Leave  Action
	IsMoved       bool
	Chain         string
	MessageRules  map[string]string
	PostbackRules map[string]string
	FallbackRule  string
}

// This gets the context data for a state
func (s *State) Data() interface{} {
	d := Get(s)
	return d
}

// Stores context for a state
func (s *State) SetData(d interface{}) {
	Set(s, d)
}

// This function forces the state change to new state
// and bypasses the Test func
func (s *State) Transit(ns string) {
	s.IsMoved = true
	s.Chain = ns
}

// Registers a regex for message transitor of this state
func (s *State) RegisterMessageRule(regex string, next string) {
	s.MessageRules[regex] = next
}

// Registers a postback transitor of this state
func (s *State) RegisterPostbackRule(payload string, next string) {
	s.PostbackRules[payload] = next
}

// Registers a fallback transitor
func (s *State) RegisterFallback(next string) {
	s.FallbackRule = next
}

// This function tests the rules defined for the state and
// returns the next state
//
// The priority order is as follows:
// 1. check message rules
// 2. check postback rules
// 3. fallback transitor
// On each step if match is found stop and return the state
func (s *State) Test(c mbotapi.Callback) (string, error) {

	if msg := c.Message; msg.Text != "" {
		for k, v := range s.MessageRules {
			match, _ := regexp.MatchString(k, msg.Text)
			if match {
				return v, nil
			}
		}
	} else if pb := c.Postback; pb != (mbotapi.InputPostback{}) {
		for k, v := range s.PostbackRules {
			if k == pb.Payload {
				return v, nil
			}
		}
	} else if s.FallbackRule != "" {
		return s.FallbackRule, nil
	}

	return "", ErrNoTransitorFound
}

// This is the function that performs action
// on entering or leaving a particular state
type Action func(c mbotapi.Callback, bot *mbotapi.BotAPI) error

// Create a new empty state
func MakeState(name string) State {
	return State{
		Name:          name,
		Enter:         nil,
		Leave:         nil,
		IsMoved:       false,
		Chain:         "",
		MessageRules:  make(map[string]string),
		PostbackRules: make(map[string]string),
		FallbackRule:  "",
	}
}
