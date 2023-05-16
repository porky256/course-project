package render

import (
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/models"
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
	td = r.addDefaultData(td, req)

	err = pageTemplate.Execute(w, td)
	if err != nil {
		log.Println("error occured while executing page template", err)
	}
}

func CreateTemplateCacheMap() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	files, err := filepath.Glob("./templates/*page.tmpl")
	if err != nil {
		return cache, fmt.Errorf("error occurred while searching for page files: %s", err)
	}

	for _, page := range files {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, fmt.Errorf("error occurred while parsing page: %s", err)
		}

		layouts, err := filepath.Glob("./templates/*layout.tmpl")
		if err != nil {
			return cache, fmt.Errorf("error occurred while searching for layout files: %s", err)
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*layout.tmpl")
			if err != nil {
				return cache, fmt.Errorf("error occurred while parsing layouts: %s", err)
			}

		}
		cache[name] = ts
	}

	return cache, nil
}

func (r *Render) addDefaultData(td *models.TemplateData, req *http.Request) *models.TemplateData {
	td.Flash = r.app.Session.PopString(req.Context(), "flash")
	td.Error = r.app.Session.PopString(req.Context(), "error")
	td.Warning = r.app.Session.PopString(req.Context(), "warning")
	td.CSRFToken = nosurf.Token(req)
	return td
}
