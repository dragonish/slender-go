package apis

import (
	"runtime"
	sIcons "slender/internal/icons"
	"slender/internal/model"
	"slender/internal/version"

	"github.com/gin-gonic/gin"
)

func about(rGroup *gin.RouterGroup) {
	rGroup.GET(model.API_ABOUT+"/info", func(ctx *gin.Context) {
		resData := model.AboutInfoData{
			Version:   version.Version,
			Commit:    version.Commit,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			BuildDate: version.BuildDate,
		}
		okWithData(ctx, resData)
	})

	rGroup.GET(model.API_ABOUT+"/icons", func(ctx *gin.Context) {
		resData := model.AboutIconsData{
			MDIVersion: sIcons.MDIVer,
			SIVersion:  sIcons.SIVer,
		}
		okWithData(ctx, resData)
	})
}
