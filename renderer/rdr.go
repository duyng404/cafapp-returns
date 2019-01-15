package renderer

import (
	"bytes"
	"cafapp-returns/logger"
	"html/template"
	"path/filepath"
)

const TEMPLATE_DIR = "./templates"
const INCLUDE_DIR = "./templates/includes/*.html"

// Rdr Renderer
type Rdr struct {
	Views map[string]view
	Fmap  template.FuncMap
}

type view struct {
	Base string
	T    *template.Template
}

// InitRdr init the renderer with a list of views and a funcmap
func InitRdr(views []string, f template.FuncMap) *Rdr {
	var r Rdr
	r.Views = make(map[string]view)
	r.Fmap = f
	for _, v := range views {
		r.addView(v)
	}
	return &r
}

func (r *Rdr) addView(name string) {
	base := filepath.Join(TEMPLATE_DIR, name+".html")
	includes, err := filepath.Glob(INCLUDE_DIR)
	if err != nil {
		logger.Fatal(err)
	}
	pages, err := filepath.Glob(filepath.Join(TEMPLATE_DIR, name, "*.html"))
	if err != nil {
		logger.Fatal(err)
	}

	for _, p := range pages {
		var tmp []string
		tmp = append(tmp, base)
		tmp = append(tmp, includes...)
		tmp = append(tmp, p)
		propername := filepath.Base(p)
		T := template.New(propername)
		T, err := T.Funcs(r.Fmap).ParseFiles(tmp...)
		if err != nil {
			logger.Fatal("failed to load template:", propername, err)
		}
		r.Views[propername] = view{
			Base: filepath.Base(base),
			T:    T,
		}
		logger.Info("Loaded template", propername)
	}

	if len(pages) == 0 {
		var tmp []string
		tmp = append(tmp, base)
		tmp = append(tmp, includes...)
		propername := filepath.Base(name + ".html")
		T := template.New(propername)
		T, err := T.Funcs(r.Fmap).ParseFiles(tmp...)
		if err != nil {
			logger.Fatal("failed to load template:", propername, err)
		}
		r.Views[propername] = view{
			Base: filepath.Base(base),
			T:    T,
		}
		logger.Info("Loaded template", propername)
	}

}

// RenderHTML render html
func (r *Rdr) RenderHTML(name string, data map[string]interface{}) (*bytes.Buffer, error) {
	logger.Info("executing templte", name)
	var buf bytes.Buffer
	// err := r.T.ExecuteTemplate(&buf, name, data)
	err := r.Views[name].T.ExecuteTemplate(&buf, r.Views[name].Base, data)
	// logger.Info(buf.String())
	return &buf, err
}
