package models

type TemplateData struct {
	StringMap  map[string]string
	IntMap     map[string]int
	Float32Map map[string]float32
	Data       map[string]interface{}
}
