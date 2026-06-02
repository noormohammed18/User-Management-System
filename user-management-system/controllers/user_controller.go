package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"user-management/middleware"
	"user-management/models"
	"user-management/repositories"
)

type PageData struct {
	Users       []models.User
	CurrentUser *models.User
	Message     string
	Error       string

	CurrentPage int
	TotalPages  int
	PrevPage    int
	NextPage    int

	TotalUsers  int
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	currentUser, err := getLoggedInUser(r)
	if err != nil || currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	page := 1

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	limit := 5

	name  := r.URL.Query().Get("name")
email := r.URL.Query().Get("email")
age   := r.URL.Query().Get("age")
role  := r.URL.Query().Get("role")

var users []models.User
if name != "" || email != "" || age != "" || role != "" {
    users, err = repositories.FilterUsers(name, email, age, role)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
} else {
    users, err = repositories.GetUsers(page, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
	

	totalUsers, err := repositories.GetUserCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalUsers + limit - 1) / limit

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	tmpl, err := template.ParseFiles("template/users.html")
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Users:       users,
		CurrentUser: currentUser,
		CurrentPage: page,
		TotalPages:  totalPages,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		TotalUsers:  totalUsers,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// helper to get logged in user from cookie
func getLoggedInUser(r *http.Request) (*models.User, error) {
	session, err := middleware.Store.Get(r, "user-session")
	if err != nil {
		return nil, err
	}
	userID, ok := session.Values["userID"].(string)
	if !ok || userID == "" {
		return nil, nil
	}
	return repositories.GetUserByID(userID)
}

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	user, err := getLoggedInUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("template/profile.html")
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, PageData{CurrentUser: user})
}


func CreateUser(w http.ResponseWriter, r *http.Request) {
	currentUser, err := getLoggedInUser(r)
	if err != nil || currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !currentUser.IsAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	age := r.FormValue("age")

	err = repositories.CreateUser(name, email, password, age)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	currentUser, err := getLoggedInUser(r)
	if err != nil || currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !currentUser.IsAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	age := r.FormValue("age")

	err = repositories.UpdateUser(id, name, email, age)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUser, err := getLoggedInUser(r)
	if err != nil || currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !currentUser.IsAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	id := r.URL.Query().Get("id")
	err = repositories.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
func ToggleAdmin(w http.ResponseWriter, r *http.Request) {
	currentUser, err := getLoggedInUser(r)
	if err != nil || currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !currentUser.IsAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	id    := r.URL.Query().Get("id")
	value := r.URL.Query().Get("value")

	err = repositories.ToggleAdmin(id, value)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func ShowEditForm(w http.ResponseWriter, r *http.Request) {
	currentUser, err := getLoggedInUser(r)
	if err != nil || currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !currentUser.IsAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	id := r.URL.Query().Get("id")
	user, err := repositories.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl, err := template.ParseFiles("template/edit.html")
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, PageData{CurrentUser: user})
}

// suppress unused import
var _ = fmt.Sprintf
