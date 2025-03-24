package pages

import (
	"slender/internal/database"
	"slender/internal/global"
	"slender/internal/model"
	"slender/internal/redirect"
	"slender/internal/validator"

	"github.com/gin-gonic/gin"
)

func logout(router *gin.Engine) {
	//? When logging out, there is no need to verify its status again
	router.GET(model.PAGE_LOGOUT, func(ctx *gin.Context) {
		accessID := validator.GetAccessID(ctx)
		if accessID != "" {
			err := database.Logout(accessID)
			if err == nil {
				ctx.SetCookie(global.Flags.GetAccessCookieName(), "", 0, model.PAGE_HOME, "", false, true)
			}
		}

		redirect.RedirectHome(ctx)
	})
}
