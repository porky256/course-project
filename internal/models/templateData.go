package models

import "github.com/porky256/course-project/internal/forms"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	Float32Map      map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
