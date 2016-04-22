package dispatcher

import "github.com/abhinavdahiya/go-messenger-bot"

type MockState struct {
	IsMoved bool
	Chain   string
}

func (m *MockState) Name() string {
	return "mock"
}

func (m *MockState) Transit(s string) {
	m.IsMoved = true
	m.Chain = s
}

func (m *MockState) Next() string {
	return m.Chain
}

func (m *MockState) Flush() {
	m.IsMoved = false
	m.Chain = ""
}

func mockEnter(state *State, c mbotapi.Callback, bot *mbotapi.BotAPI) error {
	msg := mbotapi.NewMessage("HI boss!!")
	bot.Send(c.Sender, msg, mbotapi.RegularNotif)
	return nil
}

func mockLeave(state *State, c mbotapi.Callback, bot *mbotapi.BotAPI) error {
	msg := mbotapi.NewMessage("To the next step boss!!")
	bot.Send(c.Sender, msg, mbotapi.RegularNotif)
	return nil
}

func (m *MockState) Actions() (Action, Action) {
	return mockEnter, mockLeave
}
