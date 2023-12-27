package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.davc.io/internal/models"
	"snippetbox.davc.io/internal/models/validator"

	"github.com/julienschmidt/httprouter"
)

func (self *application) about(w http.ResponseWriter, r *http.Request) {
	data := self.newTemplateData(r)
	self.render(w, http.StatusOK, "about.html", data)
}

func (self *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := self.snippets.Latest()
	if err != nil {
		self.serverError(w, err)
		return
	}

	data := self.newTemplateData(r)
	data.Snippets = snippets

	self.render(w, http.StatusOK, "home.html", data)
}

func (self *application) accountView(w http.ResponseWriter, r *http.Request) {
	userID := self.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	user, err := self.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			self.serverError(w, err)
		}
		return
	}

	data := self.newTemplateData(r)
	data.User = user

	self.render(w, http.StatusOK, "account.html", data)
}

func (self *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Retrieve named parameters from request context.
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		self.notFound(w)
		return
	}

	snippet, err := self.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			self.notFound(w)
		} else {
			self.serverError(w, err)
		}
		return
	}

	data := self.newTemplateData(r)
	data.Snippet = snippet

	self.render(w, http.StatusOK, "view.html", data)
}

// Show a snippet form
func (self *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := self.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	self.render(w, http.StatusOK, "create.html", data)
}

// Struct tags map HTML form values to struct fields.
// `form:"-"`  to ignore a field during decoding.
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// Create snippet
func (self *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := self.decodePostForm(r, &form)
	if err != nil {
		self.clientError(w, http.StatusBadRequest)
		return
	}

	// Because the Validator type is embedded by the snippetCreateForm struct,
	// we can call CheckField() directly on it.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := self.newTemplateData(r)
		data.Form = form
		self.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := self.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		self.serverError(w, err)
		return
	}

	self.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (self *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := self.newTemplateData(r)
	data.Form = userSignupForm{}
	self.render(w, http.StatusOK, "signup.html", data)
}

func (self *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := self.decodePostForm(r, &form)
	if err != nil {
		self.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := self.newTemplateData(r)
		data.Form = form
		self.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = self.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := self.newTemplateData(r)
			data.Form = form
			self.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			self.serverError(w, err)
		}

		return
	}

	self.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (self *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := self.newTemplateData(r)
	data.Form = userLoginForm{}
	self.render(w, http.StatusOK, "login.html", data)
}

func (self *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := self.decodePostForm(r, &form)
	if err != nil {
		self.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := self.newTemplateData(r)
		data.Form = form
		self.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := self.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := self.newTemplateData(r)
			data.Form = form
			self.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			self.serverError(w, err)
		}
		return
	}

	// It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user.
	// Mitigates the risk of a session fixation attacks.
	err = self.sessionManager.RenewToken(r.Context())
	if err != nil {
		self.serverError(w, err)
		return
	}

	self.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Redirect user appropriately after login
	path := self.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (self *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := self.sessionManager.RenewToken(r.Context())
	if err != nil {
		self.serverError(w, err)
		return
	}

	self.sessionManager.Remove(r.Context(), "authenticatedUserID")

	self.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type accountPasswordUpdateForm struct {
	CurrentPassword         string `form:"currentPassword"`
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

func (self *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := self.newTemplateData(r)
	data.Form = accountPasswordUpdateForm{}

	self.render(w, http.StatusOK, "password.html", data)
}

func (self *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form accountPasswordUpdateForm

	err := self.decodePostForm(r, &form)
	if err != nil {
		self.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		data := self.newTemplateData(r)
		data.Form = form
		self.render(w, http.StatusUnprocessableEntity, "password.html", data)
		return
	}

	userID := self.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = self.users.PasswordUpdate(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Current password is incorrect")
			data := self.newTemplateData(r)
			data.Form = form
			self.render(w, http.StatusUnprocessableEntity, "password.html", data)
		}
		self.serverError(w, err)
		return
	}

	self.sessionManager.Put(r.Context(), "flash", "Password successfully changed!")

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
