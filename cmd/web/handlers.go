package main

import (
	"fmt"
	"net/http"
	"strconv"

	// "text/template"
	// "time"

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

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	title := "O Saviour"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "5"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serveError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	// w.Write([]byte("Create a new snippet..."))
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
