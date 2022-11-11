package routes

type MessageLevel string

const (
	LevelPrimary MessageLevel = "primary"
	LevelSuccess MessageLevel = "success"
	LevelWarning MessageLevel = "warning"
	LevelDanger  MessageLevel = "danger"
)

// level can be only: primary, success, warning, danger. How to implement this?
type message struct {
	Text  string
	Level MessageLevel
}

var messageQueue []message

//addMessage adds message to the front of messageQueue slice
func addMessage(text string, level MessageLevel) {
	var newMessage message
	newMessage.Text = text
	newMessage.Level = level
	messageQueue = append([]message{newMessage}, messageQueue...)
}

//getMessages return all messages in messageQueue and makes it empty
func getMessages() []message {
	var tempMessages = make([]message, len(messageQueue))
	copy(tempMessages, messageQueue)
	messageQueue = []message(nil)
	return tempMessages
}
