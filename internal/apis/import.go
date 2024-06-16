package apis

import (
	"slender/internal/database"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

func apiImport(gGroup *gin.RouterGroup) {
	// import bookmarks
	gGroup.POST(model.API_IMPORT+"/bookmarks", func(ctx *gin.Context) {
		var list []model.BookmarkImportItem
		err := ctx.ShouldBindJSON(&list)
		if err == nil {
			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			count, err := database.ImportBookmarks(&list)
			if err == nil {
				created(ctx, count)
			} else {
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})
}
