package apis

import (
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

func logins(rGroup *gin.RouterGroup) {
	// get login list
	rGroup.GET(model.API_LOGINS, func(ctx *gin.Context) {
		cond := getLoginListCond(ctx)

		var body = model.LoginListData{
			List: make([]model.LoginListItem, 0),
		}

		err := database.GetLoginList(&cond, &body)
		if err == nil {
			okWithData(ctx, body)
		} else {
			internalServerError(ctx, err)
		}
	})

	// clear login log
	rGroup.DELETE(model.API_LOGINS, func(ctx *gin.Context) {
		err := database.ClearLogins()
		if err == nil {
			noContent(ctx)
		} else {
			internalServerError(ctx, err)
		}
	})

	// logout all logins
	rGroup.PATCH(model.API_LOGINS, func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				internalServerErrorWithPanic(ctx, err)
			}
		}()
		err := database.LogoutAll()
		if err == nil {
			noContent(ctx)
		} else {
			internalServerError(ctx, err)
		}
	})

	//// deprecated method
	rGroup.POST(model.API_LOGINS, func(ctx *gin.Context) {
		gone(ctx, "this method is deprecated")
	})

	// logout login
	rGroup.PATCH(model.API_LOGINS+"/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := database.Logout(id)
		if err == nil {
			noContent(ctx)
		} else {
			internalServerError(ctx, err)
		}
	})
}

func getLoginListCond(ctx *gin.Context) model.LoginListCondition {
	loginListCond := model.LoginListCondition{
		ListCondition: getListCond(ctx),
		IP:            model.MyString(ctx.Query("ip")),
		UA:            model.MyString(ctx.Query("ua")),
	}

	admin := ctx.Query("admin")
	if data.IsRouteTruthy(admin) || data.IsRouteFalsy(admin) {
		loginListCond.Admin = new(model.MyBool)
		*loginListCond.Admin = model.MyBool(data.IsRouteTruthy(admin))
	}

	return loginListCond
}
