package apis

import (
	"slender/internal/database"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

func list(rGroup *gin.RouterGroup) {
	// get base folder list
	rGroup.GET(model.API_LIST+"/folders", func(ctx *gin.Context) {
		var list = make([]model.BookmarkFolderInfo, 0)
		err := database.GetBookmarkFolderList(&list)
		if err == nil {
			okWithData(ctx, list)
		} else {
			internalServerError(ctx, err)
		}
	})
}
