package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"slender/internal/data"
	"slender/internal/logger"
)

// okWithData makes a response to successful data.
func okWithData(ctx *gin.Context, resData ...interface{}) {
	ctx.JSON(http.StatusOK, data.DataResponse(resData...))
}

// created makes a response to created data.
func created(ctx *gin.Context, resData ...interface{}) {
	ctx.JSON(http.StatusCreated, data.DataResponse(resData...))
}

// noContent makes an empty response.
// (204)
func noContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

// badRequest makes specific request error response.
func badRequest(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, errorResponse(msg))
}

// badRequestWithParse makes a response that cannot parse the request.
func badRequestWithParse(ctx *gin.Context, err error) {
	logger.WarnWithErr("parse request error", err, "path", data.ParseRequestPath(ctx.Request))
	badRequest(ctx, "Unable to parse request")
}

// unauthorized makes a unauthorized response.
func unauthorized(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusUnauthorized, errorResponse(msg))
}

// notFound makes a response with record not found.
func notFound(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusNotFound, errorResponse(msg))
}

// conflict makes a response with record conflict.
func conflict(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusConflict, errorResponse(msg))
}

// internalServerError makes an internal error response.
func internalServerError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, errorResponse(logger.ErrMsg(err)))
}

// internalServerErrorWithPanic makes an exception error response.
func internalServerErrorWithPanic(ctx *gin.Context, err interface{}) {
	msg := "unexpected error"
	if e, ok := err.(error); ok {
		logger.Err(msg, e, "path", data.ParseRequestPath(ctx.Request))
	} else {
		unknown := errors.New("nnknown, see context")
		logger.Err(msg, unknown, "context", err)
	}

	ctx.JSON(http.StatusInternalServerError, errorResponse("internal program execution error"))
}

// errorResponse generates body for error response.
func errorResponse(msg string) gin.H {
	return gin.H{
		"message": msg,
	}
}
