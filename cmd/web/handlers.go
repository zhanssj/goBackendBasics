package main

import (
	"errors"
	"fmt"
	"module1/pkg/forms"
	"module1/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.gohtml", &templateData{Snippets: s,})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.render(w, r, "notfound.page.gohtml", nil)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.gohtml", &templateData{Snippet: s,})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.gohtml", &templateData{Form: forms.New(nil)})
}

func (app *application) create(w http.ResponseWriter, r *http.Request)  {
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
		app.render(w, r, "create.page.gohtml", &templateData{Form: form})
		return
	}

	expiresNumber, err := strconv.Atoi(form.Get("expires"))
	id ,err := app.snippets.Insert(form.Get("title"), form.Get("content"), expiresNumber)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
