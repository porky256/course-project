package render

import (
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/porky256/course-project/pkg/config"
	"github.com/porky256/course-project/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Render struct {
	app *config.AppConfig
}

func NewRender(app *config.AppConfig) *Render {
	return &Render{
		app: app,
	}
}

var templateCache = make(map[string]*template.Template)

// RenderTemplateV1 deprecated
func (r *Render) RenderTemplateV1(w http.ResponseWriter, path string) {

	parsed, _ := template.ParseFiles("./templates/"+path, "./templates/base.layout.html")
	err := parsed.Execute(w, nil)
	if err != nil {
		log.Println("error parsing template:", err.Error())
	}
}

// RenderTemplateV2 deprecated
func (r *Render) RenderTemplateV2(w http.ResponseWriter, path string) {

	cached, inMap := templateCache[path]

	if !inMap {
		err := templateToCache(path)
		if err != nil {
			log.Println("error with template caching:", err.Error())
			return
		}
		cached = templateCache[path]
		log.Println("cashing:", path)
	} else {
		log.Println("using cached")
	}

	err := cached.Execute(w, nil)
	if err != nil {
		log.Println("error executing template:", err.Error())
	}
}

func (r *Render) RenderTemplateV3(w http.ResponseWriter, req *http.Request, path string, td *models.TemplateData) {
	var templateCache map[string]*template.Template
	var err error
	if !r.app.UseCache {
		templateCache, err = CreateTemplateCacheMap()
	} else {
		templateCache = r.app.TemplateCache
	}
	if err != nil {
		log.Println("error occurred while creating template cache", err)
		return
	}
	pageTemplate, ok := templateCache[path]
	if !ok {
		log.Println("asked page not found:", path)
		return
	}
	td = AddDefaultData(td, req)

	err = pageTemplate.Execute(w, td)
	if err != nil {
		log.Println("error occured while executing page template", err)
	}
}

func CreateTemplateCacheMap() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	files, err := filepath.Glob("./templates/*page.html")
	if err != nil {
		return cache, fmt.Errorf("error occurred while searching for page files: %s", err)
	}

	for _, page := range files {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, fmt.Errorf("error occurred while parsing page: %s", err)
		}

		layouts, err := filepath.Glob("./templates/*layout.html")
		if err != nil {
			return cache, fmt.Errorf("error occurred while searching for layout files: %s", err)
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*layout.html")
			if err != nil {
				return cache, fmt.Errorf("error occurred while parsing layouts: %s", err)
			}

		}
		cache[name] = ts
	}

	return cache, nil
}

func templateToCache(path string) error {
	_, inMap := templateCache[path]
	if inMap {
		return fmt.Errorf("cached template already exists: %s", path)
	}
	parsed, err := template.ParseFiles("./templates/"+path, "./templates/base.layout.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %s", err)
	}
	templateCache[path] = parsed
	return nil
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}
