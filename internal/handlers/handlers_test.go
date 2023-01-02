package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	Key   string
	Value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majsu", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"contact", "/contacts", "GET", []postData{}, http.StatusOK},

	{"post-search-availability", "/search-availability", "POST", []postData{
		{Key: "start_date", Value: "2020-01-01"},
		{Key: "end_date", Value: "2020-01-02"},
	}, http.StatusOK},

	{"post-search-availability-json", "/search-availability-json", "POST", []postData{
		{Key: "start_date", Value: "2020-01-01"},
		{Key: "end_date", Value: "2020-01-02"},
	}, http.StatusOK},

	{"post_make_reservation", "/make-reservation", "POST", []postData{
		{Key: "first_name", Value: "Kishor"},
		{Key: "last_name", Value: "Padmanabhan"},
		{Key: "email", Value: "kishor@123.com"},
		{Key: "phone", Value: "9947766456"},
	}, http.StatusOK},
}

func TestHandler(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
		if e.method == "POST" {
			values := url.Values{}

			for _, x := range e.params {

				values.Add(x.Key, x.Value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
