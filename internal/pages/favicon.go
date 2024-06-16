package pages

import (
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

// favicon registers favicon routing.
func favicon(router *gin.Engine) {
	router.StaticFile("/favicon.ico", model.FAVICON_FILE)
}
