package auth

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	SessionUsername string = "username"
	SessionRole     string = "role"
)

var SessionStore *session.Store

func CreateSession() {
	SessionStore = session.New(session.Config{})
}
