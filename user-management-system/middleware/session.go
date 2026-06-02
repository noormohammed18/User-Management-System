package middleware

import "github.com/gorilla/sessions"

// Store is the global session store.
// Change the secret key to a long random string in production.
var Store = sessions.NewCookieStore([]byte("super-secret-key-change-in-production"))

func init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 day in seconds
		HttpOnly: true,  // not accessible via JS
	}
}