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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func login(router *gin.Engine) {
	header := "Welcome to Slender"
	action := "/login"
	username := "slender"

	router.GET(model.PAGE_LOGIN, accessBypasser, func(ctx *gin.Context) {
		localizer := local.New(ctx)

		title, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Login",
				Other: "Login",
			},
		})
		ok, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "OK",
				Other: "OK",
			},
		})
		placeholder, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AccessPlaceholder",
				Other: "Enter access password",
			},
		})

		ctx.HTML(http.StatusOK, "login.go.tmpl", gin.H{
			"Title":       title + " - " + global.Config.Title,
			"OK":          ok,
			"Placeholder": placeholder,
			"Header":      header,
			"Action":      action,
			"Username":    username,
		})
	})

	router.POST(model.PAGE_LOGIN, accessBypasser, func(ctx *gin.Context) {
		pwd := ctx.PostForm("password")

		localizer := local.New(ctx)

		title := localizer.Message("Login") + " - " + global.Config.Title
		ok := localizer.Message("OK")
		placeholder := localizer.Message("AccessPlaceholder")

		emptyTip, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "PwdEmpty",
				Other: "Password is empty",
			},
		})
		incorrectTip, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "PwdIncorrect",
				Other: "Incorrect password",
			},
		})

		if pwd == "" {
			ctx.HTML(http.StatusUnauthorized, "login.go.tmpl", gin.H{
				"Title":       title,
				"OK":          ok,
				"Placeholder": placeholder,
				"Tip":         emptyTip,
				"Header":      header,
				"Action":      action,
				"Username":    username,
			})
			ctx.Abort()
			return
		}

		readIp := ip.GetRealIP(ctx)
		ua := ctx.GetHeader("User-Agent")
		now := time.Now()
		expires := now.AddDate(0, 0, int(global.Flags.TokenAge))

		logger.Info("logging in", "login_time", now, "ip", readIp, "user_agent", ua)

		if pwd == global.Flags.AccessPassword {
			uid, err := uuid.NewV4()
			if err == nil {
				requestID := uid.String()

				claims := data.ClaimsGenerator(requestID, global.Flags.AccessToken, now, expires)
				jwt := data.JWTGenerator(global.Flags.Secret, claims)

				err := database.AddLogin(requestID, now, readIp, ua, false)
				if err != nil {
					//* It will not affect the successful status of login.
					logger.Warn("this login was not recorded in the database")
				}

				ctx.SetCookie(model.COOKIE_ACCESS_PREFIX+global.Flags.GetPortStr(), jwt, global.Flags.GetTokenAgeSeconds(), model.PAGE_HOME, "", false, true)

				redirect.RedirectHome(ctx)
			} else {
				ctx.HTML(http.StatusInternalServerError, "login.go.tmpl", gin.H{
					"Title":       title,
					"OK":          ok,
					"Placeholder": placeholder,
					"Tip":         logger.ErrMsg(logger.Err("unable to generate id", err)),
					"Header":      header,
					"Action":      action,
					"Username":    username,
				})
			}
		} else {
			ctx.HTML(http.StatusUnauthorized, "login.go.tmpl", gin.H{
				"Title":       title,
				"OK":          ok,
				"Placeholder": placeholder,
				"Tip":         incorrectTip,
				"Header":      header,
				"Action":      action,
				"Username":    username,
			})
		}

		ctx.Abort()
	})
}
