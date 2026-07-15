package main

import (
	"net/http"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	err := app.html.render(w, http.StatusOK, nil, "base", "pages/home.tmpl")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) gopher(w http.ResponseWriter, r *http.Request) {
	width := 100
	err := app.html.render(w, http.StatusOK, width, "partial:image:gopher")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// Define a user type. The fields need to be exported so that we can reference
// them in our HTML templates.
type user struct {
	Name     string
	Email    string
	IsGopher bool
}

// Create a hardcoded list of users.
var users = []user{
	{"Alice Madsen", "alice.madsen@example.com", true},
	{"Theo Thatcher", "theo.thatcher@example.com", true},
	{"Maxwell Albright", "maxwell.albright@example.com", false},
	{"Ruby Thompson", "ruby.thompson@example.com", false},
	{"Leona Rowan", "leona.rowan@example.com", false},
	{"Alicia Lennox", "alicia.lennox@example.com", true},
	{"Ruben Mason", "ruben.mason@example.com", false},
	{"Leo Reynolds", "leo.reynolds@example.com", false},
	{"Max Lester", "max.lester@example.com", true},
	{"Theodore Allister", "theodore.allister@example.com", false},
}

func (app *application) listUsers(w http.ResponseWriter, r *http.Request) {
	err := app.html.render(w, http.StatusOK, users, "base", "pages/users.tmpl")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) searchUsers(w http.ResponseWriter, r *http.Request) {
	// Filter down the list of users to find ones that match the query.
	query := r.FormValue("query")

	matches := users
	if query != "" {
		matches = []user{}
		for _, u := range users {
			if strings.Contains(u.Name, query) || strings.Contains(u.Email, query) {
				matches = append(matches, u)
			}
		}
	}

	// Render the base template by default, but render the users:rows template instead if it's a HTMX request (a partial).
	template := "base"
	if isHTMXRequest(r) {
		template = "users:rows"
	}

	err := app.html.render(w, http.StatusOK, matches, template, "pages/users.tmpl")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
