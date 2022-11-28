package session

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

const SessionUser string = "username"

var SessionStore *session.Store

func CreateSession() {
	SessionStore = session.New(session.Config{})
}
