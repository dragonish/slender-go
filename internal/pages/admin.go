package pages

import (
	"net/http"
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/global"
	"slender/internal/ip"
	"slender/internal/local"
	"slender/internal/logger"
	"slender/internal/model"
	"slender/internal/redirect"
	"slender/internal/validator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func admin(router *gin.Engine) {
	header := "Administrator Mode"
	action := "/admin?redirect="
	username := "admin"

	router.GET(model.PAGE_ADMIN, adminBypasser, func(ctx *gin.Context) {
		rURL := ctx.DefaultQuery("redirect", model.PAGE_MANAGER)

		localizer := local.New(ctx)

		title, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Admin",
				Other: "Admin",
			},
		})
		placeholder, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AdminPlaceholder",
				Other: "Enter admin password",
			},
		})

		ctx.HTML(http.StatusOK, "login.go.tmpl", gin.H{
			"Title":       title + " - " + global.Config.Title,
			"OK":          localizer.Message("OK"),
			"Placeholder": placeholder,
			"Header":      header,
			"Action":      action + rURL,
			"Username":    username,
		})
	})

	// login admin
	router.POST(model.PAGE_ADMIN, adminBypasser, func(ctx *gin.Context) {
		rURL := ctx.DefaultQuery("redirect", model.PAGE_MANAGER)

		pwd := ctx.PostForm("password")

		localizer := local.New(ctx)

		title := localizer.Message("Admin") + " - " + global.Config.Title
		ok := localizer.Message("OK")
		placeholder := localizer.Message("AdminPlaceholder")

		if pwd == "" {
			ctx.HTML(http.StatusUnauthorized, "login.go.tmpl", gin.H{
				"Title":       title,
				"OK":          ok,
				"Placeholder": placeholder,
				"Tip":         localizer.Message("PwdEmpty"),
				"Header":      header,
				"Action":      action + rURL,
				"Username":    username,
			})
			ctx.Abort()
			return
		}

		readIp := ip.GetRealIP(ctx)
		ua := ctx.GetHeader("User-Agent")
		now := time.Now()
		expires := now.AddDate(0, 0, int(global.Flags.TokenAge))

		logger.Info("admin logining", "login_time", now, "ip", readIp, "user_agent", ua)

		if pwd == global.Flags.AdminPassword {
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

				redirect.Redirect(ctx, rURL)
			} else {
				ctx.HTML(http.StatusInternalServerError, "login.go.tmpl", gin.H{
					"Title":       title,
					"OK":          ok,
					"Placeholder": placeholder,
					"Tip":         logger.ErrMsg(logger.Err("unable to generate id", err)),
					"Header":      header,
					"Action":      action + rURL,
					"Username":    username,
				})
			}
		} else {
			ctx.HTML(http.StatusUnauthorized, "login.go.tmpl", gin.H{
				"Title":       title,
				"OK":          ok,
				"Placeholder": placeholder,
				"Tip":         localizer.Message("PwdIncorrect"),
				"Header":      header,
				"Action":      action + rURL,
				"Username":    username,
			})
		}
	})

	// logout admin
	//? When logging out, there is no need to verify its status again
	router.GET(model.PAGE_ADMIN+"/logout", func(ctx *gin.Context) {
		adminID := validator.GetAdminID(ctx)
		if adminID != "" {
			err := database.Logout(adminID)
			if err == nil {
				ctx.SetCookie(global.Flags.GetAdminCookieName(), "", -1, model.PAGE_HOME, "", false, true)
			}
		}

		redirect.RedirectHome(ctx)
	})
}
