package forms

type errors map[string][]string

// add errors
func(e errors) Add(field, message string){
	e[field] = append(e[field], message)
}

// get errors
func(e errors) Get(field string) string{
	es:= e[field]

	if(len(es)==0){
		return "";
	}
	return es[0]
}