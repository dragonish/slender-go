package apis

import (
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/global"
	"slender/internal/ip"
	"slender/internal/logger"
	"slender/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

// admin handles routes that do not require permissions.
func admin(rGroup *gin.RouterGroup) {
	// request administrator status
	rGroup.POST(model.API_ADMIN, adminBypasser, func(ctx *gin.Context) {
		var body model.RequestAdminPostBoby
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			if body.Password == "" {
				badRequest(ctx, "password is empty")
				ctx.Abort()
				return
			}

			readIp := ip.GetRealIP(ctx)
			ua := ctx.GetHeader("User-Agent")
			now := time.Now()
			expires := now.AddDate(0, 0, int(global.Flags.TokenAge))

			logger.Info("admin logining", "login_time", now, "ip", readIp, "user_agent", ua)

			if body.Password == global.Flags.AdminPassword {
				uid, err := uuid.NewV4()
				if err == nil {
					requestID := uid.String()

					claims := data.ClaimsGenerator(requestID, global.Flags.AdminToken, now, expires)
					jwt := data.JWTGenerator(global.Flags.Secret, claims)

					err := database.AddLogin(requestID, now, readIp, ua, true, global.Flags.TokenAge)
					if err != nil {
						//* It will not affect the successful status of login.
						logger.Warn("this login was not recorded in the database")
					}

					ctx.SetCookie(global.Flags.GetAdminCookieName(), jwt, global.Flags.GetTokenAgeSeconds(), model.PAGE_HOME, "", false, true)

					okWithData(ctx, jwt)
				} else {
					internalServerError(ctx, logger.Err("unable to generate id", err))
				}
			} else {
				unauthorized(ctx, "incorrect password")
			}
		} else {
			badRequestWithParse(ctx, err)
		}

		ctx.Abort()
	})

	// Get current admin id
	rGroup.GET(model.API_ADMIN, adminHandler, func(ctx *gin.Context) {
		adminId := ctx.GetString(model.CONTEXT_ADMIN_ID)
		if adminId == "" {
			internalServerError(ctx, logger.NewErr("unable to get admin id"))
		} else {
			okWithData(ctx, adminId)
		}
	})
}
