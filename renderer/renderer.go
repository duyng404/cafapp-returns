package renderer

import (
	"cafapp-returns/logger"
	"html/template"
	"os"
	"path/filepath"
)

const TEMPLATE_FOLDER = "templates"

type View struct {
	Name string
	T    *template.Template
}

var views map[string]View

func addView(name string) {
	v := View{
		Name: name,
		T:    template.New(name),
	}
	pattern := filepath.Join(TEMPLATE_FOLDER, name+"-*.tmpl")
	v.T = template.New(name)
	_, err := v.T.ParseGlob(pattern)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(v.T.DefinedTemplates())
	views[name] = v
}

func init() {
	views = make(map[string]View)
	addView("landing")
}

func RenderView(name string) {
	err := views[name].T.ExecuteTemplate(os.Stdout, "landing-base", nil)
	if err != nil {
		logger.Error(err)
	}
}
