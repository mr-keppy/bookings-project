package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/mr-keppy/bookings/pkg/config"
	"github.com/mr-keppy/bookings/pkg/models"
)

//RenderTemplate renders using html template
func RenderTemplateTest(w http.ResponseWriter, templ string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + templ,"./templates/base.layout.tmpl")

	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}


var app *config.AppConfig

func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData{
	td.CSRFToken = nosurf.Token(r)
	return td
}


func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData){

	var tc map[string] *template.Template

	if(app.UseCache){
	tc= app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]

	if !ok{
		log.Fatal("Error while writing template")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)

	if err!=nil{
		log.Println((err))
	}
	_, err = buf.WriteTo(w)

	if err!=nil{
		log.Println(err)
	}


}

func CreateTemplateCache()(map[string]*template.Template, error){
	myCache := map[string]*template.Template{}

	// get all of the files name from *page.temp from template
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if(err !=nil){
		return myCache, err
	}

	//range through all files
	for _, page:= range pages{
		name := filepath.Base(page)
		ts, err:= template.New(name).ParseFiles((page))
		if(err !=nil){
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if(err !=nil){
			return myCache, err
		}

		if(len(matches)>0){
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
		}

		if(err !=nil){
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