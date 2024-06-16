package validator

import (
	"slender/internal/data"
	"slender/internal/global"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

// AccessValidator returns true when access or admin verification passed.
func AccessValidator(ctx *gin.Context) (valid bool) {
	valid = AdminValidator(ctx)
	if valid {
		return
	}

	port := global.Flags.GetPortStr()

	accessCookie, _ := ctx.Cookie(model.COOKIE_ACCESS_PREFIX + port)
	if accessCookie != "" {
		accessJWT, _ := data.ParseJWT(global.Flags.Secret, accessCookie)
		if accessJWT.Token == global.Flags.AccessToken {
			valid = true
			ctx.Set(model.CONTEXT_IDENTITY, "access")
			return
		}
	}

	return
}

// AdminValidator returns true when admin verification passed.
func AdminValidator(ctx *gin.Context) (valid bool) {
	port := global.Flags.GetPortStr()

	adminCookie, _ := ctx.Cookie(model.COOKIE_ADMIN_PREFIX + port)
	if adminCookie != "" {
		adminJWT, _ := data.ParseJWT(global.Flags.Secret, adminCookie)
		if adminJWT.Token == global.Flags.AdminToken {
			valid = true
			ctx.Set(model.CONTEXT_IDENTITY, "admin")
		}
	}

	return
}
