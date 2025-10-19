package apis

import (
	"slender/internal/database"
	"slender/internal/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func search_engines(rGroup *gin.RouterGroup) {
	// add search engine
	rGroup.POST(model.API_SEARCH_ENGINES, func(ctx *gin.Context) {
		var body model.SearchEnginePostBody
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			if body.Name == "" {
				badRequest(ctx, "search engine name is empty")
				return
			}

			if body.URL == "" {
				badRequest(ctx, "search engine url is empty")
				return
			}

			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			searchEngineID, err := database.AddSearchEngine(&body)
			if err == nil {
				created(ctx, searchEngineID)
			} else {
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}

		ctx.Abort()
	})

	// get search engine list
	// optional order values: created-time | modified-time | weight.
	rGroup.GET(model.API_SEARCH_ENGINES, func(ctx *gin.Context) {
		cond := getSearchEngineListCond(ctx)

		var body = model.SearchEngineListData{
			List: make([]model.SearchEngineBaseData, 0),
		}

		err := database.GetSearchEngines(&cond, &body)
		if err == nil {
			okWithData(ctx, body)
		} else {
			internalServerError(ctx, err)
		}
	})

	// handle search engine in batches
	// action "delete" | "setWeight" | "icnWeight"
	rGroup.PATCH(model.API_SEARCH_ENGINES, func(ctx *gin.Context) {
		var body model.BatchPatchBody
		err := ctx.ShouldBindJSON(&body)
		if err == nil {
			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			err = database.SearchEngineBatchHandler(&body)
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

	// get search engine
	rGroup.GET(model.API_SEARCH_ENGINES+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			var body = model.SearchEngineBaseData{}

			err := database.GetSearchEngine(id, &body)
			switch err {
			case nil:
				okWithData(ctx, body)
			case model.ErrNotExist:
				notFound(ctx, "search engine does not exist")
			default:
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})

	// update search engine
	rGroup.PATCH(model.API_SEARCH_ENGINES+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			var body model.SearchEnginePatchBody
			err := ctx.ShouldBindJSON(&body)
			if err == nil {
				defer func() {
					if err := recover(); err != nil {
						internalServerErrorWithPanic(ctx, err)
					}
				}()
				err = database.UpdateSearchEngine(id, &body)
				switch err {
				case nil:
					noContent(ctx)
				case model.ErrNotExist:
					notFound(ctx, "search does not exist")
				case model.ErrDoNothing:
					badRequest(ctx, "invalid reques data")
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

	// delete search engine
	rGroup.DELETE(model.API_SEARCH_ENGINES+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			defer func() {
				if err := recover(); err != nil {
					internalServerErrorWithPanic(ctx, err)
				}
			}()
			err := database.DeleteSearchEngine(id)
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

func getSearchEngineListCond(ctx *gin.Context) model.SearchEngineListCondition {
	searchEngineListCond := model.SearchEngineListCondition{
		ListCondition: getListCond(ctx),
		Name:          model.MyString(ctx.Query("name")),
		URL:           model.MyString(ctx.Query("url")),
	}

	method := strings.ToLower(ctx.Query("method"))
	if method == "get" || method == "post" {
		searchEngineListCond.Method = new(model.MyString)
		*searchEngineListCond.Method = model.MyString(method)
	}

	return searchEngineListCond
}
