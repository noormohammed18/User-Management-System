package controllers

import (
	"html/template"
	"net/http"
	"fmt"
	"user-management/middleware"
	"user-management/repositories"
)

type LoginPageData struct {
	Error   string
	Success string
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("template/login.html")
		if err != nil {
			http.Error(w, "Template not found: "+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, LoginPageData{})
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := repositories.GetUserByEmailAndPassword(email, password)
		if err != nil || user == nil {
			tmpl, _ := template.ParseFiles("template/login.html")
			tmpl.Execute(w, LoginPageData{Error: "Invalid email or password"})
			return
		}

		// Save userID in session
		session, _ := middleware.Store.Get(r, "user-session")
		session.Values["userID"] = fmt.Sprintf("%d", user.ID)
		session.Save(r, w)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func ShowRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("template/register.html")
		if err != nil {
			http.Error(w, "Template not found: "+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, LoginPageData{})
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		age := r.FormValue("age")

		err := repositories.CreateUser(name, email, password, age)
		if err != nil {
			tmpl, _ := template.ParseFiles("template/register.html")
			tmpl.Execute(w, LoginPageData{Error: "Registration failed: " + err.Error()})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the session
	session, _ := middleware.Store.Get(r, "user-session")
	session.Values["userID"] = nil
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}