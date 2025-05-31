package pages

import (
	"html/template"
	"net/http"
	"slender/internal/global"
	"slender/internal/local"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// home makes the homepage content response.
func home(router *gin.Engine) {
	// homepage
	router.GET(model.PAGE_HOME, accessHandler, func(ctx *gin.Context) {
		localizer := local.New(ctx)

		ungrouped, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Ungrouped",
				Other: "Ungrouped",
			},
		})
		latest, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Latest",
				Other: "Latest",
			},
		})
		hot, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Hot",
				Other: "Hot",
			},
		})
		clearSearchTip, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ClearSearchTip",
				Other: "Click this or press ESC to clear",
			},
		})
		inHomeSearch, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "InHomeSearch",
				Other: "In Home Search",
			},
		})
		useInHomeSearch, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "UseInHomeSearch",
				Other: "Enable in home search",
			},
		})
		folders, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Folders",
				Other: "Folders",
			},
		})
		admin := localizer.Message("Admin")
		user, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "User",
				Other: "User",
			},
		})
		privacy, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Privacy",
				Other: "Privacy",
			},
		})
		quit, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Quit",
				Other: "Quit",
			},
		})
		manager, _ := localizer.L.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Manager",
				Other: "Manager",
			},
		})

		dynamic := model.PageDynamicURL{}
		dynamic.Parse(ctx.Request)

		searchModule := generateSearchModule(clearSearchTip, inHomeSearch, useInHomeSearch)

		identity := ctx.GetString(model.CONTEXT_IDENTITY)
		bookmarks, sidebar := generateBookmarks(&dynamic, identity == "admin", ungrouped, latest, hot)

		ctx.HTML(http.StatusOK, "home.go.tmpl", gin.H{
			"Title":           global.Config.Title,
			"ShowSidebar":     global.Config.ShowSidebar,
			"ShowSearchInput": global.Config.ShowSearchInput,
			"ShowScrollTop":   global.Config.ShowScrollTop,
			"CustomFooter":    template.HTML(global.Config.CustomFooter),

			"FoldersText": folders,
			"AdminText":   admin,
			"UserText":    user,
			"PrivacyText": privacy,
			"QuitText":    quit,
			"ManagerText": manager,

			"SearchModule": searchModule,
			"Bookmarks":    bookmarks,
			"Sidebar":      sidebar,

			"Identity": identity,
		})
	})
}
