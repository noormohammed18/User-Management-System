package middleware

import "net/http"

// AuthMiddleware protects routes by checking for a valid session.
// If the session is missing or userID is not set, redirects to /login.
func LoginAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "user-session")
		if err != nil || session.Values["userID"] == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func LogoutAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "user-session")
		if err == nil && session.Values["userID"] != nil {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}