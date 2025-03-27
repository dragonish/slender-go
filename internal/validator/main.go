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

	accessCookie, _ := ctx.Cookie(global.Flags.GetAccessCookieName())
	if accessCookie != "" {
		accessJWT, _ := data.ParseJWT(global.Flags.Secret, accessCookie)
		if accessJWT.Token == global.Flags.AccessToken && global.Flags.IsLogined(accessJWT.Subject) {
			valid = true
			ctx.Set(model.CONTEXT_IDENTITY, "access")
			return
		}
	}

	return
}

// AdminValidator returns true when admin verification passed.
func AdminValidator(ctx *gin.Context) (valid bool) {
	adminCookie, _ := ctx.Cookie(global.Flags.GetAdminCookieName())
	if adminCookie != "" {
		adminJWT, _ := data.ParseJWT(global.Flags.Secret, adminCookie)
		if adminJWT.Token == global.Flags.AdminToken && global.Flags.IsLogined(adminJWT.Subject) {
			valid = true
			ctx.Set(model.CONTEXT_IDENTITY, "admin")
			ctx.Set(model.CONTEXT_ADMIN_ID, adminJWT.Subject)
		}
	}

	return
}

// GetAccessID returns the access id.
func GetAccessID(ctx *gin.Context) string {
	accessCookie, _ := ctx.Cookie(global.Flags.GetAccessCookieName())
	if accessCookie != "" {
		accessJWT, _ := data.ParseJWT(global.Flags.Secret, accessCookie)
		return accessJWT.Subject
	}

	return ""
}

// GetAdminID returns the admin id.
func GetAdminID(ctx *gin.Context) string {
	adminCookie, _ := ctx.Cookie(global.Flags.GetAdminCookieName())
	if adminCookie != "" {
		adminJWT, _ := data.ParseJWT(global.Flags.Secret, adminCookie)
		return adminJWT.Subject
	}

	return ""
}
