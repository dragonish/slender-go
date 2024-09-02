package pages

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"slender/internal/icons"
	"slender/internal/logger"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/psanford/memfs"
)

var mFs *memfs.FS

// MEM_ICONS_DIR defines icons dir in memory.
const MEM_ICONS_DIR = "assets/icons"

func Init() {
	logger.Debug("page init")
	mFs = memfs.New()
	err := mFs.MkdirAll(MEM_ICONS_DIR, 0777)
	if err != nil {
		logger.Panic("unable to create icons folder in memory", err)
	}

	for siName, siContent := range icons.SimpleIcons {
		file := filepath.ToSlash(filepath.Join(MEM_ICONS_DIR, "si-"+siName+".svg"))
		svg := `<svg role="img" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="` + siContent + `"></path></svg>`

		err := mFs.WriteFile(file, []byte(svg), 0755)
		if err != nil {
			logger.Err("unable to write icon file in memory", err, "group", "simple-icons", "name", siName)
		}
	}

	for mdiName, mdiContent := range icons.MaterialDesignIcons {
		file := filepath.ToSlash(filepath.Join(MEM_ICONS_DIR, "mdi-"+mdiName+".svg"))
		svg := `<svg role="img" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="` + mdiContent + `"></path></svg>`

		err := mFs.WriteFile(file, []byte(svg), 0755)
		if err != nil {
			logger.Err("unable to write icon file in memory", err, "group", "material-design-icons", "name", mdiName)
		}
	}
}

func Pages(router *gin.Engine) {
	Init()

	htmlTemplates := []string{model.TEMPLATES_PATH + "/home.go.tmpl", model.TEMPLATES_PATH + "/login.go.tmpl", model.MANAGER_PATH + "/manager.html"}

	router.LoadHTMLFiles(htmlTemplates...)
	iconsFs, _ := fs.Sub(mFs, MEM_ICONS_DIR)
	uploadsFs := gin.Dir(model.UPLOAD_FILES_PATH, false)

	router.Static("/assets/js", model.JS_PATH)
	router.Static("/assets/css", model.CSS_PATH)
	router.Static("/manager-assets", model.MANAGER_PATH+"/manager-assets")

	router.GET("/assets/icons/*filepath", func(ctx *gin.Context) {
		filepath := ctx.Param("filepath")
		setCacheHeader(ctx)
		ctx.FileFromFS(filepath, http.FS(iconsFs))
	})

	router.GET("/assets/uploads/*filepath", func(ctx *gin.Context) {
		filepath := ctx.Param("filepath")
		setCacheHeader(ctx)
		ctx.FileFromFS(filepath, uploadsFs)
	})

	home(router)
	login(router)
	logout(router)
	admin(router)
	manager(router)
	favicon(router)
}
