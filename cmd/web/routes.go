package main

import (
	"net/http"

	"snippetbox.davc.io/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (self *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		self.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// For testing
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// For session management, create a new middleware chain.
	dynamic := alice.New(self.sessionManager.LoadAndSave, noSurf, self.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(self.home))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(self.about))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(self.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(self.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(self.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(self.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(self.userLoginPost))

	protected := dynamic.Append(self.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(self.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(self.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(self.userLogoutPost))
	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(self.accountView))
	router.Handler(http.MethodGet, "/account/password/update", protected.ThenFunc(self.accountPasswordUpdate))
	router.Handler(http.MethodPost, "/account/password/update", protected.ThenFunc(self.accountPasswordUpdatePost))

	// Middleware chaining
	// recoverPanic -> logRequest -> secureHeaders -> app
	return alice.New(self.recoverPanic, self.logRequest, secureHeaders).Then(router)
}
