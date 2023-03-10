package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Required(fields ...string){
	for _, field := range fields{
			value:= f.Get(field)

			if(strings.TrimSpace(value)==""){
				f.Errors.Add(field, "this field cannot be blank")
			}
	}
}

func New(data url.Values) *Form{
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string, r *http.Request) bool{
	x:= f.Get(field)
	if x==""{
		return false
	}
	return true
}

// return true if no errors
func (f *Form) Valid() bool{
	return len(f.Errors)==0
}

// check min length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x:= f.Get(field)

	if len(x)<length{
		f.Errors.Add(field,fmt.Sprintf("this field must be at least %d characters long",length))
		return false
	}
	return true
}

// email validator
func (f *Form) IsEmail(field string){
	if !govalidator.IsEmail(f.Get(field)){
		f.Errors.Add(field,"invalid email address")
	}
}