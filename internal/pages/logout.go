package pages

import (
	"slender/internal/global"
	"slender/internal/model"
	"slender/internal/redirect"

	"github.com/gin-gonic/gin"
)

func logout(router *gin.Engine) {
	router.GET(model.PAGE_LOGOUT, accessHandler, func(ctx *gin.Context) {
		ctx.SetCookie(model.COOKIE_ACCESS_PREFIX+global.Flags.GetPortStr(), "", 0, model.PAGE_HOME, "", false, true)
		redirect.RedirectHome(ctx)
	})
}
