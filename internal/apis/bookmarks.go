package apis

import (
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func bookmarks(rGroup *gin.RouterGroup) {
	// add bookmark
	rGroup.POST(model.API_BOOKMARKS, func(ctx *gin.Context) {
		var body model.BookmarkPostBody
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			if body.URL == "" {
				badRequest(ctx, "bookmark url is empty")
				return
			}

			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			bookmarkID, err := database.AddBookmark(&body)
			if err == nil {
				created(ctx, bookmarkID)
			} else {
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}

		ctx.Abort()
	})

	// get bookmark list
	// optional order values: created-time | modified-time | visits | folder-weight | weight.
	rGroup.GET(model.API_BOOKMARKS, func(ctx *gin.Context) {
		cond := getBookmarkListCond(ctx)

		var body = model.BookmarkListData{
			List: make([]model.BookmarkListItem, 0),
		}

		err := database.GetBookmarkList(&cond, &body)
		if err == nil {
			okWithData(ctx, body)
		} else {
			internalServerError(ctx, err)
		}
	})

	// handle bookmark in batches
	// action: "delete" | "setPrivacy" | "setWeight" | "incWeight" | "clearVisits" | "setFolder" | "setHideInOther"
	rGroup.PATCH(model.API_BOOKMARKS, func(ctx *gin.Context) {
		var body model.BatchPatchBody
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			err = database.BookmarkBatchHandler(&body)
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

	// get bookmark
	rGroup.GET(model.API_BOOKMARKS+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			var body = model.BookmarkResData{
				Files: make([]model.FileBaseData, 0),
			}

			err := database.GetBookmark(id, &body)
			switch err {
			case nil:
				okWithData(ctx, body)
			case model.ErrNotExist:
				notFound(ctx, "bookmark does not exist")
			default:
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})

	// update bookmark
	rGroup.PATCH(model.API_BOOKMARKS+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			var body model.BookmarkPatchBody
			err := ctx.ShouldBindJSON(&body)
			if err == nil {
				defer func() {
					if err := recover(); err != nil {
						internalServerErrorWithPanic(ctx, err)
					}
				}()
				err = database.UpdateBookmark(id, &body)
				switch err {
				case nil:
					noContent(ctx)
				case model.ErrNotExist:
					notFound(ctx, "bookmark does not exist")
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

	// delete bookmark
	rGroup.DELETE(model.API_BOOKMARKS+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			err := database.DeleteBookmark(id)
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

func getBookmarkListCond(ctx *gin.Context) model.BookmarkListCondition {
	bookmarkListCond := model.BookmarkListCondition{
		ListCondition: getListCond(ctx),
		Name:          model.MyString(ctx.Query("name")),
		Des:           model.MyString(ctx.Query("description")),
		URL:           model.MyString(ctx.Query("url")),
	}

	privacy := ctx.Query("privacy")
	if data.IsRouteTruthy(privacy) || data.IsRouteFalsy(privacy) {
		bookmarkListCond.Privacy = new(model.MyBool)
		*bookmarkListCond.Privacy = model.MyBool(data.IsRouteTruthy(privacy))
	}

	folder := ctx.Query("folder")
	if f, err := strconv.ParseInt(folder, 10, 64); err == nil {
		if f >= 0 {
			bookmarkListCond.Folder = new(model.NullInt64)
			*bookmarkListCond.Folder = model.NullInt64(f)
		}
	}

	hideInOther := ctx.Query("hide-in-other")
	if data.IsRouteTruthy(hideInOther) || data.IsRouteFalsy(hideInOther) {
		bookmarkListCond.HideInOther = new(model.MyBool)
		*bookmarkListCond.HideInOther = model.MyBool(data.IsRouteTruthy(hideInOther))
	}

	return bookmarkListCond
}
