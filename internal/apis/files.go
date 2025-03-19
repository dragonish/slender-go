package apis

import (
	"os"
	"path"
	"slender/internal/data"
	"slender/internal/database"
	"slender/internal/logger"
	"slender/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func files(rGroup *gin.RouterGroup) {
	// upload file
	rGroup.POST(model.API_FILES, func(ctx *gin.Context) {
		if !data.IsPathExists(model.UPLOAD_FILES_PATH + "/") {
			logger.Info("create upload directory")
			err := os.MkdirAll(model.UPLOAD_FILES_PATH, 0777)
			if err != nil {
				internalServerError(ctx, logger.Err("create upload directory error", err, "dir", model.UPLOAD_FILES_PATH))
			}
		}

		file, _ := ctx.FormFile("file")
		timeStr := strconv.FormatInt(time.Now().UnixMilli(), 10) + "-" + data.BoundaryGenerator(6)
		filename := timeStr + path.Ext(file.Filename)

		log := logger.New("filename", file.Filename)
		savePath := model.UPLOAD_FILES_PATH + "/" + filename
		err := ctx.SaveUploadedFile(file, savePath)
		if err == nil {

			//* record file.
			//! there is no associated bookmark at this time.
			fileID, err := database.AddFile(filename, 0)
			if err == nil {
				webPath := "/assets/uploads/" + filename
				log.Info("upload file", "web_path", webPath, "file_id", fileID)

				resData := model.FileBaseData{
					ID:   model.MyInt64(fileID),
					Path: model.MyString(webPath),
				}
				okWithData(ctx, resData)
			} else {
				//! unable to record, delete file.
				data.DeleteFile(savePath)
				internalServerError(ctx, log.Err("unable to record file", err))
			}

		} else {
			internalServerError(ctx, log.Err("save uploaded file error", err))
		}
	})

	// get file list
	rGroup.GET(model.API_FILES, func(ctx *gin.Context) {
		cond := getFileListCond(ctx)

		var body = model.FileListData{
			List: make([]model.FileListItem, 0),
		}

		err := database.GetFileList(&cond, &body)
		if err == nil {
			okWithData(ctx, body)
		} else {
			internalServerError(ctx, err)
		}
	})

	// remove all unused files
	rGroup.DELETE(model.API_FILES, func(ctx *gin.Context) {
		err := database.RemoveUnusedFiles()
		if err == nil {
			noContent(ctx)
		} else {
			internalServerError(ctx, err)
		}
	})

	// delete file
	rGroup.DELETE(model.API_FILES+"/:id", func(ctx *gin.Context) {
		id, err := parseIDParam(ctx, "id")
		if err == nil {
			force := data.IsRouteTruthy(ctx.Query("force"))

			err := database.DeleteFile(id, force)
			if err == nil {
				noContent(ctx)
			} else if err == model.ErrDoNothing {
				conflict(ctx, "the file cannot be deleted")
			} else {
				internalServerError(ctx, err)
			}
		} else {
			badRequestWithParse(ctx, err)
		}
	})
}

func getFileListCond(ctx *gin.Context) model.FileListCondition {
	fileListCond := model.FileListCondition{
		ListCondition: getListCond(ctx),
		Path:          model.MyString(ctx.Query("path")),
	}

	use := ctx.Query("use")
	if data.IsRouteTruthy(use) || data.IsRouteFalsy(use) {
		fileListCond.Use = new(model.MyBool)
		*fileListCond.Use = model.MyBool(data.IsRouteTruthy(use))
	}

	return fileListCond
}
