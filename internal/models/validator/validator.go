package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Map of validation errors for form fields.
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Check if valid email address.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (self *Validator) AddFieldError(key, message string) {
	if self.FieldErrors == nil {
		self.FieldErrors = make(map[string]string)
	}

	_, exists := self.FieldErrors[key]
	if !exists {
		self.FieldErrors[key] = message
	}
}

func (self *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		self.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
