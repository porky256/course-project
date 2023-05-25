package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/porky256/course-project/internal/models"
	"html/template"
	"log"
)

type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	IsProduction  bool
	Session       *scs.SessionManager
	RootPath      string
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	DateLayout    string
	MailChan      chan models.MailData
}
