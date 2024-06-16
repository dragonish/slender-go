package pages

import (
	"net/http"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

func manager(router *gin.Engine) {
	router.GET(model.PAGE_MANAGER, adminHandler, func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "manager.html", nil)
	})
}
