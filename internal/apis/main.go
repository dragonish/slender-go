package apis

import (
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

// Apis registers all API Routing.
func Apis(router *gin.Engine) {
	//? for frontend development.
	devGroup := router.Group(model.PREFIX_V1)
	admin(devGroup)
	about(devGroup)

	accessGroup := router.Group(model.PREFIX_V1, accessHandler)
	access(accessGroup)

	adminGroup := router.Group(model.PREFIX_V1, adminHandler)
	apiImport(adminGroup)
	config(adminGroup)
	bookmarks(adminGroup)
	files(adminGroup)
	folders(adminGroup)
	search_engines(adminGroup)
	icons(adminGroup)
	list(adminGroup)
	logins(adminGroup)
}
