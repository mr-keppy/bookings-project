package models

import "github.com/mr-keppy/bookings/internal/forms"

//holds data set
type TemplateData struct{
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float32
	Data map[string]interface{}
	CSRFToken string
	Flash string
	Warning string
	Error string
	Form *forms.Form
	IsAuthenticated int
}
