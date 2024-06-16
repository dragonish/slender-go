package redirect

import (
	"net/http"
	"slender/internal/data"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

// Redirect makes a response with Redirect url.
func Redirect(ctx *gin.Context, url string) {
	ctx.Header("Location", url)
	ctx.JSON(http.StatusSeeOther, data.DataResponse(url))
}

// RedirectHome makes a response with redirect homepage.
func RedirectHome(ctx *gin.Context) {
	Redirect(ctx, model.PAGE_HOME)
}

// RedirectLogin makes a response with redirect login page.
func RedirectLogin(ctx *gin.Context) {
	Redirect(ctx, model.PAGE_LOGIN)
}

// RedirectAdmin makes a response with redirect admin page.
func RedirectAdmin(ctx *gin.Context) {
	Redirect(ctx, model.PAGE_ADMIN)
}
