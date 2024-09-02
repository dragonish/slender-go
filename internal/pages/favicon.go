package pages

import (
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

// favicon registers favicon routing.
func favicon(router *gin.Engine) {
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		setCacheHeader(ctx)
		ctx.File(model.FAVICON_FILE)
	})
}
