package apis

import (
	conf "slender/internal/config"
	"slender/internal/global"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

func config(rGroup *gin.RouterGroup) {
	// get config
	rGroup.GET(model.API_CONFIG, func(ctx *gin.Context) {
		okWithData(ctx, global.Config)
	})

	// update config
	rGroup.PUT(model.API_CONFIG, func(ctx *gin.Context) {
		var body model.UserConfig
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			conf.Update(body)
			noContent(ctx)
		} else {
			badRequestWithParse(ctx, err)
		}
	})

	// patch update config
	rGroup.PATCH(model.API_CONFIG, func(ctx *gin.Context) {
		var body model.ConfigPatchBody
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			conf.PatchUpdate(body)
			noContent(ctx)
		} else {
			badRequestWithParse(ctx, err)
		}
	})
}
