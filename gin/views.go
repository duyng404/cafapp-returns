package gin

import (
	"cafapp-returns/logger"
	"html/template"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

const TEMPLATE_DIR = "./templates"
const INCLUDE_DIR = "./templates/includes/*.html"

func initViews(router *gin.Engine) {
	templateList := []string{
		"landing",
		"userdash",
		"order",
		"404",
	}
	router.HTMLRender = loadTemplates(templateList)
}

func loadTemplates(list []string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	f := template.FuncMap{
		"formatMoney": formatMoney,
		"rawHTML":     rawHTML,
	}

	for _, name := range list {
		base := filepath.Join(TEMPLATE_DIR, name+".html")
		includes, err := filepath.Glob(INCLUDE_DIR)
		if err != nil {
			logger.Fatal(err)
		}
		views, err := filepath.Glob(filepath.Join(TEMPLATE_DIR, name, "*.html"))
		if err != nil {
			logger.Fatal(err)
		}

		for _, v := range views {
			var tmp []string
			tmp = append(tmp, base)
			tmp = append(tmp, includes...)
			tmp = append(tmp, v)
			r.AddFromFilesFuncs(filepath.Base(v), f, tmp...)
			logger.Info("Loaded template", filepath.Base(v))
		}

		if len(views) == 0 {
			var tmp []string
			tmp = append(tmp, base)
			tmp = append(tmp, includes...)
			r.AddFromFilesFuncs(filepath.Base(name+".html"), f, tmp...)
			logger.Info("Loaded template", filepath.Base(name+".html"))
		}
	}

	return r
}

func formatMoney(a uint) string {
	l := a / 100
	r := a % 100
	ls := strconv.FormatUint(uint64(l), 10)
	rs := strconv.FormatUint(uint64(r), 10)
	if r < 10 {
		rs = "0" + rs
	}
	return "$" + ls + "." + rs
}

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}
