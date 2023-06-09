package render

import (
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/helpers"
	"github.com/porky256/course-project/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"humanDate":  humanDate,
	"formatTime": formatTime,
	"makeRange":  makeRange,
}

type Render struct {
	app *config.AppConfig
}

const dateLayout = "2006-01-02"

func humanDate(t time.Time) string {
	return t.Format(dateLayout)
}

func formatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

func makeRange(start, end, step int) []int {
	var ans []int
	for i := start; i <= end; i += step {
		ans = append(ans, i)
	}
	return ans
}

// NewRender creates new render entity
func NewRender(app *config.AppConfig) *Render {
	return &Render{
		app: app,
	}
}

// Template renders a page based on template
func (r *Render) Template(w http.ResponseWriter, req *http.Request, path string, td *models.TemplateData) error {
	var templateCache map[string]*template.Template
	var err error
	if !r.app.UseCache {
		templateCache, err = CreateTemplateCacheMap(r.app)
	} else {
		templateCache = r.app.TemplateCache
	}
	if err != nil {
		helpers.ServerError(w, fmt.Errorf("error occurred while creating template cache: %e", err))
		return fmt.Errorf("error occurred while creating template cache: %e", err)
	}
	pageTemplate, ok := templateCache[path]
	if !ok {
		helpers.ServerError(w, fmt.Errorf("asked page not found: %s", path))
		return fmt.Errorf("asked page not found: %s", path)
	}
	td = r.addDefaultData(td, req)

	err = pageTemplate.Execute(w, td)
	if err != nil {
		helpers.ServerError(w, fmt.Errorf("error occurred while executing page template: %s", err))
		return fmt.Errorf("error occurred while executing page template: %s", err)
	}
	return nil
}

// CreateTemplateCacheMap creates cache for page templates
func CreateTemplateCacheMap(app *config.AppConfig) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	files, err := filepath.Glob(app.RootPath + "/static/templates/*page.tmpl")
	if err != nil {
		app.ErrorLog.Println("error occurred while searching for page files:", err)
		return cache, fmt.Errorf("error occurred while searching for page files: %s", err)
	}

	for _, page := range files {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			app.ErrorLog.Printf("error occurred while parsing page %s:%e \n", name, err)
			return cache, fmt.Errorf("error occurred while parsing page %s:%e", name, err)
		}

		layouts, err := filepath.Glob(app.RootPath + "/static/templates/*layout.tmpl")
		if err != nil {
			app.ErrorLog.Println("error occurred while searching for layout files:", err)
			return cache, fmt.Errorf("error occurred while searching for layout files: %s", err)
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(app.RootPath + "/static/templates/*layout.tmpl")
			if err != nil {
				app.ErrorLog.Println("error occurred while parsing layouts:", err)
				return cache, fmt.Errorf("error occurred while parsing layouts: %s", err)
			}

		}
		cache[name] = ts
	}

	return cache, nil
}

// addDefaultData adds default data
func (r *Render) addDefaultData(td *models.TemplateData, req *http.Request) *models.TemplateData {
	td.Flash = r.app.Session.PopString(req.Context(), "flash")
	td.Error = r.app.Session.PopString(req.Context(), "error")
	td.Warning = r.app.Session.PopString(req.Context(), "warning")
	td.CSRFToken = nosurf.Token(req)
	//if r.app.Session.Exists(req.Context(), "user_id") {
	//	td.IsAuthenticated = 1
	//}
	td.IsAuthenticated = 1
	return td
}
