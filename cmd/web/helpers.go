package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// Common dynamic data.
func (self *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           self.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: self.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (self *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = self.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

// Render templates from cache
func (self *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := self.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		self.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	// Write template to buffer, instead of straight to http.ResponseWriter.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		self.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

// Request contexts, are a way to pass data alongside a HTTP request
// as it is processed by handlers or middleware. This data could be a
// user ID, a CSRF token, a web token, whether a user is logged in
// or not — something typically derived from logic that
// you don't want to repeat over-and-over again in every handler.
func (self *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (self *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	self.errorLog.Output(2, trace)

	if self.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (self *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (self *application) notFound(w http.ResponseWriter) {
	self.clientError(w, http.StatusNotFound)
}
