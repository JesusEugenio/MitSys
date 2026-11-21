package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var SessionStore *sessions.CookieStore

func InitStore(sessionKey string) {
	SessionStore = sessions.NewCookieStore([]byte(sessionKey))

	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}
