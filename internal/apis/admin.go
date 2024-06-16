package apis

import (
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/global"
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

			ip := ctx.ClientIP()
			ua := ctx.GetHeader("User-Agent")
			now := time.Now()
			expires := now.AddDate(0, 0, int(global.Flags.TokenAge))

			logger.Info("admin logining", "login_time", now, "ip", ip, "user_agent", ua)

			if body.Password == global.Flags.AdminPassword {
				uid, err := uuid.NewV4()
				if err == nil {
					requestID := uid.String()

					claims := data.ClaimsGenerator(requestID, global.Flags.AdminToken, now, expires)
					jwt := data.JWTGenerator(global.Flags.Secret, claims)

					err := database.AddLogin(requestID, now, ip, ua, true)
					if err != nil {
						//* It will not affect the successful status of login.
						logger.Warn("this login was not recorded in the database")
					}

					ctx.SetCookie(model.COOKIE_ADMIN_PREFIX+global.Flags.GetPortStr(), jwt, global.Flags.GetTokenAgeSeconds(), model.PAGE_HOME, "", false, true)

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
}
