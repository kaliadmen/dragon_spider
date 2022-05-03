package dragonSpider

import (
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

func (ds *DragonSpider) Validator(data url.Values) *Validation {
	return &Validation{
		Data:   data,
		Errors: make(map[string]string),
	}
}

func (v *Validation) InsertError(key, msg string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = msg
	}
}

func (v *Validation) Has(r *http.Request, field string) bool {
	f := r.Form.Get(field)
	if f == "" {
		return false
	}
	return true
}

func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.InsertError(field, "This field is required")
		}
	}
}

func (v *Validation) Validate(ok bool, key, msg string) {
	if !ok {
		v.InsertError(key, msg)
	}
}

func (v *Validation) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validation) IsValidEmail(field, val string) {
	if !govalidator.IsEmail(val) {
		v.InsertError(field, "Invalid email address")
	}
}

func (v *Validation) IsInt(field, val string) {
	_, err := strconv.Atoi(val)

	if err != nil {
		v.InsertError(field, "This field has to be an integer")
	}
}

func (v *Validation) IsFloat(field, val string) {
	_, err := strconv.ParseFloat(val, 64)

	if err != nil {
		v.InsertError(field, "This field has to be a floating point number(decimal)")
	}
}

func (v *Validation) IsISODate(field, val string) {
	_, err := time.Parse("2006-01-02", val)
	if err != nil {
		v.InsertError(field, "The date must be in YYYY-MM-DD format")
	}
}

func (v *Validation) NoSpaces(field, val string) {
	if govalidator.HasWhitespace(val) {
		v.InsertError(field, "spaces are not allowed in this field")
	}
}
