package main

import (
	"bytes"
	"fmt"
	"github.com/ganeshrahul23/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(data, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) addDefaultData(data *templateData, r *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}
	data.CSRFToken = nosurf.Token(r)
	data.AuthenticatedUser = app.authenticatedUser(r)
	data.CurrentYear = time.Now().Year()
	data.Flash = app.session.PopString(r, "flash")
	return data
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
