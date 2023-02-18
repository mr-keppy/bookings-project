package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/models"
)

var pathToTemplates = "./templates"
// RenderTemplate renders using html template
func RenderTemplateTest(w http.ResponseWriter, templ string) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+templ, "./templates/base.layout.tmpl")

	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

var functions = template.FuncMap{
	"humanDate": HumanDate,
}

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

// return formated date
func HumanDate(t time.Time) string {
	return t.Format("2006-01-01")
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if(app.Session.Exists(r.Context(),"user_id")){
		td.IsAuthenticated = 1
	}else{
		td.IsAuthenticated = 0
	}
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		//log.Fatal("Could not get template from template cache")
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files name from *page.temp from template
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	//range through all files
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles((page))
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		}

		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}
	return myCache, nil
}

/*func createTemplateCache_1(t string) error{

	templates := []string{
		fmt.Sprintf("./templates/%s",t),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)

	if err!=nil{
		return err
	}
	//tc[t] = tmpl

	return nil
}*/
