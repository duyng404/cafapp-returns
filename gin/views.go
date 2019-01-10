package gin

import (
	"cafapp-returns/logger"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

const TEMPLATE_DIR = "./templates"
const INCLUDE_DIR = "./templates/includes/*.html"

func initViews(router *gin.Engine) {
	templateList := []string{
		"landing",
	}
	router.HTMLRender = loadTemplates(templateList)
}

func loadTemplates(list []string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	for _, name := range list {
		base := filepath.Join(TEMPLATE_DIR, name+".html")
		includes, err := filepath.Glob(INCLUDE_DIR)
		if err != nil {
			logger.Panic(err)
		}
		views, err := filepath.Glob(filepath.Join(TEMPLATE_DIR, name, "*.html"))
		if err != nil {
			logger.Panic(err)
		}

		for _, v := range views {
			var tmp []string
			tmp = append(tmp, base)
			tmp = append(tmp, includes...)
			tmp = append(tmp, v)
			t := r.AddFromFiles(filepath.Base(v), tmp...)
			logger.Info("Loaded template", filepath.Base(v), t.DefinedTemplates())
		}
	}

	return r
}
