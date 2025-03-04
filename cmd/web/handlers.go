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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%s\n", snippet)
	// }

	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// if err != nil || id < 1 {
	// 	app.notFound(w)
	// 	return
	// }

	// s, err := app.snippets.Get(id)

	// if err == models.ErrNoRecord {
	// 	app.notFound(w)
	// 	return
	// } else if err != nil {
	// 	app.serveError(w, err)
	// 	return
	// }

	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	s, err := app.snippets.Latest()
	if err != nil {
		app.serveError(w, err)
		return
	}

	// Create an instance of a templateData struct holding the slice of
	// Snippts.
	// data := &templateData{Snippets: s}

	// files := []string{
	// 	"./ui/html/show.page.tmpl",
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serveError(w, err)
	// 	return
	// }

	// // w.Header().Set("Content-Type", "text/html; charsest=utf-8")

	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serveError(w, err)
	// }

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
	//  flash := app.session.PopString(r, "flash" )

	// app.render(w, r, "show.page.tmpl", &templateData{Flash: flash, Snippet: s})
	app.render(w, r, "show.page.tmpl", &templateData{Snippet: s})

	// data := &templateData{Snippet: s, CurrentYear: time.Now().Year()}
	// files := []string{
	// 	"./ui/html/show.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// ts, ok := app.templateCache["show.page.tmpl"]
	// if !ok {
	// 	app.serveError(w, err)
	// 	return
	// }

	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serveError(w, err)
	// }

	// fmt.Fprintf(w, "Display a specific snippet with ID %d... ", id)
	// fmt.Fprintf(w, "%s", s)
}

// func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
// 	// if r.Method != "POST" {
// 	// 	w.Header().Set("Allow", "POST")
// 	// 	app.clientError(w, http.StatusMethodNotAllowed)
// 	// 	return
// 	// }

// 	title := "O Saviour"
// 	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
// 	expires := "5"

// 	id, err := app.snippets.Insert(title, content, expires)
// 	if err != nil {
// 		app.serveError(w, err)
// 		return
// 	}

// 	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

// 	// w.Write([]byte("Create a new snippet..."))
// }

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

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires := r.PostForm.Get("expires")

	// errors := make(map[string]string)

	// if strings.TrimSpace(title) == "" {
	// 	errors["title"] = "This field cannot be blank"
	// } else if utf8.RuneCountInString(title) > 100 {
	// 	errors["title"] = "This field is too long (maximum is 100 characters)"
	// } else if utf8.RuneCountInString(title) < 3 {
	// 	errors["title"] = "This field is too small (minimum is 3 characters)"
	// }

	// if strings.TrimSpace(content) == "" {
	// 	errors["content"] = "This field cannot be blank"
	// }

	// if strings.TrimSpace(expires) == "" {
	// 	errors["expires"] = "This field cannot be blank"
	// } else if expires != "365" && expires != "7" && expires != "1" {
	// 	errors["expires"] = "This field is invalid"
	// }

	// if len(errors) > 0 {
	// 	app.render(w, r, "create.page.tmpl", &templateData{
	// 		FormErrors: errors,
	// 		FormData:   r.PostForm,
	// 	})
	// 	return
	// }

	// // Insert paresed form into the databae
	// id, err := app.snippets.Insert(title, content, expires)
	// if err != nil {
	// 	app.serveError(w, err)
	// 	return
	// }

	// http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
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
		app.render(w, r, "signup.page.tmpl", &templateData{Form:form})
		return
	}else  if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "flash", "your signup was successful. Please log in.")
	fmt.Fprintln(w, "Creating a new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "authenticate and login user...")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Logging out the user....")
}
