package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// Form custom form struct that embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes new form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid checks if there is any errors in form
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks if form fields is in post and not empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		if strings.TrimSpace(f.Get(field)) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength check field for minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d symbols long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "This field is not a valid email")
		return false
	}
	return true
}
