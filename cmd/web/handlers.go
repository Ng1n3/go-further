package main

import (
	"fmt"
	"net/http"
	"strconv"

	// "text/template"
	// "time"

	"github.com/Ng1n3/go-further/pkg/forms"
	"github.com/Ng1n3/go-further/pkg/models"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{Snippet: s})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serveError(w, err)
		return
	}

	// Use the Put() mtehod to add a string value and corresponding key to the session data. NOte that if there's no existing session for the current user then a new, empty, session for them will authomatically be created by the session middleware.

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

/*
Note: When we check the length of the title field, we’re using the
utf8.RuneCount() function — not the builtin len() function. This is
because we want to count the number of characters in the title rather
than the number of bytes. To illustrate the difference, the string "Zoë"
has 3 characters but a length of 4 bytes because of the umlauted ë
character.
*/

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupuser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRx)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "flash", "your signup was successful. Please log in.")
	fmt.Fprintln(w, "Creating a new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	//add id of current user to session
	app.session.Put(r, "userID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")

	app.session.Put(r, "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
