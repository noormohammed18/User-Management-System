package routes

import (
	"net/http"
	"user-management/controllers"
	"user-management/middleware"
)

func SetupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	// Public routes (no auth required)
	http.HandleFunc("/login", middleware.LogoutAuthMiddleware(controllers.ShowLogin))
	http.HandleFunc("/register", middleware.LogoutAuthMiddleware(controllers.ShowRegister))
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/forgot-password", controllers.ForgotPassword)
	http.HandleFunc("/verify-otp", controllers.VerifyOTP)
	http.HandleFunc("/reset-password", controllers.ResetPassword)

	// Protected routes (login required)
	http.HandleFunc("/profile", middleware.LoginAuthMiddleware(controllers.ShowProfile))
	http.HandleFunc("/users", middleware.LoginAuthMiddleware(controllers.GetUsers))

	// Protected admin routes
	http.HandleFunc("/create-user", middleware.LoginAuthMiddleware(controllers.CreateUser))
	http.HandleFunc("/edit-user", middleware.LoginAuthMiddleware(controllers.ShowEditForm))
	http.HandleFunc("/update-user", middleware.LoginAuthMiddleware(controllers.UpdateUser))
	http.HandleFunc("/delete-user", middleware.LoginAuthMiddleware(controllers.DeleteUser))
	http.HandleFunc("/toggle-admin", middleware.LoginAuthMiddleware(controllers.ToggleAdmin))
}
