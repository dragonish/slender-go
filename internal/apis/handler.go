package apis

import (
	"slender/internal/data"
	"slender/internal/global"
	"slender/internal/model"
	"slender/internal/validator"
	"strconv"

	"github.com/gin-gonic/gin"
)

// accessHandler defines access validation middleware.
func accessHandler(ctx *gin.Context) {
	if global.Flags.AccessPassword == "" {
		return
	}

	if !validator.AccessValidator(ctx) {
		unauthorized(ctx, "invalid certificate")
		ctx.Abort()
	}
}

// adminHandler defines admin validation middleware.
func adminHandler(ctx *gin.Context) {
	if global.Flags.AdminPassword == "" {
		return
	}

	if !validator.AdminValidator(ctx) {
		unauthorized(ctx, "invalid certificate")
		ctx.Abort()
	}
}

// adminBypasser defines admin bypasser.
func adminBypasser(ctx *gin.Context) {
	if global.Flags.AdminPassword == "" {
		//* bypass without an admin password.
		noContent(ctx)
		ctx.Abort()
		return
	}

	if validator.AdminValidator(ctx) {
		//* bypass when there is a valid certificate.
		noContent(ctx)
		ctx.Abort()
		return
	}
}

// parseIDParam returns int64 type id param.
func parseIDParam(ctx *gin.Context, paramName string) (int64, error) {
	idParam := ctx.Param(paramName)
	id, err := strconv.ParseInt(idParam, 10, 64)
	return id, err
}

// getListCond returns list condition.
func getListCond(ctx *gin.Context) model.ListCondition {
	order := ctx.DefaultQuery("order", "")
	size := ctx.DefaultQuery("size", "25")
	page := ctx.DefaultQuery("page", "1")

	s := data.GetSizeFromStr(size, 25, 1, 100)
	p := data.GetPageFromStr(page)

	listCond := model.ListCondition{
		Order: order,
		Page:  p,
		Size:  s,
	}

	return listCond
}
