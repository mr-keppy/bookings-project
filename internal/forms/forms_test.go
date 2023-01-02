package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)


func TestForm_Valid(t *testing.T) {

	r:= httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if(!isValid){
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {

	r:= httptest.NewRequest("POST","/whatever",nil)
	form := New(r.PostForm)

	form.Required("a","b","c")
	if(form.Valid()){
		t.Error("form shows valid when required fields are empty")
	}

	postedData:= url.Values{}

	postedData.Add("a","a")
	postedData.Add("b","a")
	postedData.Add("c","a")

	r, _ = http.NewRequest("POST","/watever",nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a","b","c")

	if !form.Valid(){
		t.Error("show does not have required when")
	}

}

func TestForm_Has(t *testing.T){
	r:= httptest.NewRequest("POST","/watever",nil)

	form:= New(r.PostForm)

	has:= form.Has("watever",r)

	if(has){
		t.Error("forms shows has field when it does not")
	}

	postedData:= url.Values{}

	postedData.Add("a","a")

	form = New(postedData)

	has = form.Has("a",r)

	if !has {
		t.Error("show forms does not have a filed when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	r:= httptest.NewRequest("POST","/watever",nil)

	form := New(r.PostForm)

	form.MinLength("x",10,r)

	if(form.Valid()){
		t.Error("forms shows min length for non existing field")
	}

	isError:= form.Errors.Get("x")
	if(isError==""){
		t.Error("should have an error but did not get one")
	}

	postedData:= url.Values{}
	postedData.Add("a","a")

	form = New(postedData)

	form.MinLength("a",5,r)

	if(form.Valid()){
		t.Error("forms shows should show min length error")
	}

	postedData = url.Values{}
	postedData.Add("some_values","abcd1234")
	form = New(postedData)

	form.MinLength("some_values",1,r)

	if(!form.Valid()){
		t.Error("forms shows should show min length not met")
	}
	isError= form.Errors.Get("some_values")
	if(isError!=""){
		t.Error("should not have an error but got one")
	}

}
func TestForm_IsEmail(t *testing.T) {
	r:= httptest.NewRequest("POST","/watever",nil)

	form := New(r.PostForm)

	form.IsEmail("x")

	if(form.Valid()){
		t.Error("forms shows min length for non existing field")
	}

	postedData := url.Values{}
	postedData.Add("email","abcd1234")
	form = New(postedData)

	form.IsEmail("email")

	if(form.Valid()){
		t.Error("forms should show invalid email error")
	}

	postedData = url.Values{}
	postedData.Add("new_email","abcd@1234.com")
	form = New(postedData)

	form.IsEmail("new_email")

	if(!form.Valid()){
		t.Error("forms should show valid email")
	}
}