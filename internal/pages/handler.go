package pages

import (
	"slender/internal/global"
	"slender/internal/model"
	"slender/internal/redirect"
	"slender/internal/validator"

	"github.com/gin-gonic/gin"
)

// accessHandler defines access validation middleware.
func accessHandler(ctx *gin.Context) {
	if global.Flags.AccessPassword == "" {
		return
	}

	if !validator.AccessValidator(ctx) {
		redirect.RedirectLogin(ctx)
		ctx.Abort()
	}
}

// adminHandler defines admin validation middleware.
func adminHandler(ctx *gin.Context) {
	if global.Flags.AdminPassword == "" {
		return
	}

	if !validator.AdminValidator(ctx) {
		redirect.RedirectAdmin(ctx)
		ctx.Abort()
	}
}

// accessBypasser defines access bypasser.
func accessBypasser(ctx *gin.Context) {
	if global.Flags.AccessPassword == "" {
		//* redirect homepage without an access password.
		redirect.RedirectHome(ctx)
		ctx.Abort()
		return
	}

	if validator.AccessValidator(ctx) {
		//* redirect homepage when there is a valid certificate.
		redirect.RedirectHome(ctx)
		ctx.Abort()
		return
	}
}

// adminBypasser defines admin bypasser.
func adminBypasser(ctx *gin.Context) {
	rURL := ctx.DefaultQuery("redirect", model.PAGE_MANAGER)

	if global.Flags.AdminPassword == "" {
		//* redirect url without an admin password.
		redirect.Redirect(ctx, rURL)
		ctx.Abort()
		return
	}

	if validator.AdminValidator(ctx) {
		//* redirect url when there is a valid certificate.
		redirect.Redirect(ctx, rURL)
		ctx.Abort()
		return
	}
}
