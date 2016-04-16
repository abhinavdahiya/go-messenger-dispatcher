package dispatcher

import (
	"regexp"

	"github.com/abhinavdahiya/go-messenger-bot"
)

type MockState struct {
	ID           string
	Enter, Leave Action

	IsMoved bool
	Chain   string
}

func (m *MockState) Name() string {
	return m.ID
}

func (m *MockState) Transit(s string) {
	m.IsMoved = true
	m.Chain = s
}

func (m *MockState) Next() string {
	return m.Chain
}

func (m *MockState) Transitor(c mbotapi.Callback) {
	if m.IsMoved {
		return
	}

	if msg := c.Message; msg.Text != "" {
		if match, _ := regexp.MatchString("(?i)*hi*", msg.Text); match {
			m.IsMoved = true
			m.Chain = "Hi"
			return
		}
	}
}

func (m *MockState) Flush() {
	m.IsMoved = false
	m.Chain = ""
}

func mockEnter(c mbotapi.Callback, bot *mbotapi.BotAPI) error {
	msg := mbotapi.NewMessage("HI boss!!")
	bot.Send(c.Sender, msg, mbotapi.RegularNotif)
	return nil
}

func mockLeave(c mbotapi.Callback, bot *mbotapi.BotAPI) error {
	msg := mbotapi.NewMessage("To the next step boss!!")
	bot.Send(c.Sender, msg, mbotapi.RegularNotif)
	return nil
}

func (m *MockState) Actions() (Action, Action) {
	return mockEnter, mockLeave
}
