package controllers

import (
	"html/template"
	"net/http"
	"user-management/repositories"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("template/forgot_password.html")
		tmpl.Execute(w, nil)
		return
	}

	email := r.FormValue("email")

	otp := repositories.GenerateOTP()

	err := repositories.SaveOTP(email, otp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = SendOTP(email, otp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(
		w,
		r,
		"/reset-password?email="+email,
		http.StatusSeeOther,
	)
}
