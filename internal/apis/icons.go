package apis

import (
	sIcons "slender/internal/icons"
	"slender/internal/model"
	"sort"

	"github.com/gin-gonic/gin"
)

func icons(rGroup *gin.RouterGroup) {
	// get material-design-icons list
	rGroup.GET(model.API_ICONS, func(ctx *gin.Context) {
		list := make([]string, 0)
		for key := range sIcons.MaterialDesignIcons {
			list = append(list, "mdi-"+key)
		}
		for key := range sIcons.SimpleIcons {
			list = append(list, "si-"+key)
		}
		// sort
		sort.Strings(list)
		okWithData(ctx, list)
	})
}
