package apis

import (
	"slender/internal/database"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

// access handles routes that only requires access premission.
func access(rGroup *gin.RouterGroup) {
	// increase bookmark visits
	rGroup.POST(model.API_BOOKMARKS+"/:id/visits", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			err := database.IncreaseBookmarkVisits(id)
			if err == nil {
				noContent(ctx)
			} else {
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})
}
