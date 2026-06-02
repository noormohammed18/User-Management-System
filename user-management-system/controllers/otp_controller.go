package controllers

import (
	"html/template"
	"net/http"
	"user-management/repositories"
)

type OTPPageData struct {
	Email string
	Error string
}

type ResetPasswordData struct {
	Email   string
	Error   string
	Success string
}

// VerifyOTP handles both GET and POST /verify-otp
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/verify_otp.html")
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), 500)
		return
	}

	if r.Method == http.MethodGet {
		email := r.URL.Query().Get("email")
		tmpl.Execute(w, OTPPageData{Email: email})
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		otp := r.FormValue("otp")

		valid, err := repositories.VerifyOTP(email, otp)
		if err != nil {
			tmpl.Execute(w, OTPPageData{Email: email, Error: "Something went wrong. Please try again."})
			return
		}

		if !valid {
			tmpl.Execute(w, OTPPageData{Email: email, Error: "Invalid or expired OTP. Please try again."})
			return
		}

		_ = repositories.DeleteOTP(email)
		http.Redirect(w, r, "/reset-password?email="+email, http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// ResetPassword handles GET and POST /reset-password
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/reset_password.html")
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), 500)
		return
	}

	if r.Method == http.MethodGet {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Redirect(w, r, "/forgot-password", http.StatusSeeOther)
			return
		}
		tmpl.Execute(w, ResetPasswordData{Email: email})
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			tmpl.Execute(w, ResetPasswordData{Email: email, Error: "Passwords do not match."})
			return
		}

		if len(password) < 6 {
			tmpl.Execute(w, ResetPasswordData{Email: email, Error: "Password must be at least 6 characters."})
			return
		}

		err := repositories.UpdatePasswordByEmail(email, password)
		if err != nil {
			tmpl.Execute(w, ResetPasswordData{Email: email, Error: "Failed to update password. Please try again."})
			return
		}

		tmpl.Execute(w, ResetPasswordData{Success: "Your password has been reset successfully!"})
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}