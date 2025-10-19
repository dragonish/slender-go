package apis

import (
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/model"

	"github.com/gin-gonic/gin"
)

func folders(rGroup *gin.RouterGroup) {
	// add folder
	rGroup.POST(model.API_FOLDERS, func(ctx *gin.Context) {
		var body model.FolderPostBody
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			if body.Name == "" {
				badRequest(ctx, "folder name is empty")
				return
			}

			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			folderID, err := database.AddFolder(&body)
			if err == nil {
				created(ctx, folderID)
			} else {
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}

		ctx.Abort()
	})

	// get folder list
	// optional order values: created-time | modified-time | bookmark-total | weight.
	rGroup.GET(model.API_FOLDERS, func(ctx *gin.Context) {
		cond := getFolderListCond(ctx)

		var body = model.FolderListData{
			List: make([]model.FolderListItem, 0),
		}

		err := database.GetFolderList(&cond, &body)
		if err == nil {
			okWithData(ctx, body)
		} else {
			internalServerError(ctx, err)
		}
	})

	// handle folder in batches
	// action: "delete" | "setLarge" | "setPrivacy" | "setWeight" | "incWeight"
	rGroup.PATCH(model.API_FOLDERS, func(ctx *gin.Context) {
		var body model.BatchPatchBody
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			err = database.FolderBatchHandler(&body)
			switch err {
			case nil:
				noContent(ctx)
			case model.ErrDoNothing:
				badRequest(ctx, "unable to recognize action")
			case model.ErrQueryParamMissing:
				badRequest(ctx, "invalid payload")
			default:
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}

		ctx.Abort()
	})

	// get folder
	rGroup.GET(model.API_FOLDERS+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			var body model.FolderBaseData
			err := database.GetFolder(id, &body)
			switch err {
			case nil:
				okWithData(ctx, body)
			case model.ErrNotExist:
				notFound(ctx, "folder does not exist")
			default:
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})

	// update folder
	rGroup.PATCH(model.API_FOLDERS+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			var body model.FolderPatchBody
			err := ctx.ShouldBindJSON(&body)
			if err == nil {
				err = database.UpdateFolder(id, &body)
				switch err {
				case nil:
					noContent(ctx)
				case model.ErrNotExist:
					notFound(ctx, "folder does not exist")
				case model.ErrDoNothing:
					badRequest(ctx, "invalid request data")
				default:
					internalServerError(ctx, err)
				}
			} else {
				badRequestWithParse(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})

	// delete folder
	rGroup.DELETE(model.API_FOLDERS+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			err = database.DeleteFolder(id)
			if err == nil {
				noContent(ctx)
			} else {
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})
}

func getFolderListCond(ctx *gin.Context) model.FolderListCondition {
	folderListCond := model.FolderListCondition{
		ListCondition: getListCond(ctx),
		Name:          model.MyString(ctx.Query("name")),
		Des:           model.MyString(ctx.Query("description")),
	}

	privacy := ctx.Query("privacy")
	if data.IsRouteTruthy(privacy) || data.IsRouteFalsy(privacy) {
		folderListCond.Privacy = new(model.MyBool)
		*folderListCond.Privacy = model.MyBool(data.IsRouteTruthy(privacy))
	}

	return folderListCond
}
