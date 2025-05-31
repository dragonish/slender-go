package data

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ParseRequestPath returns the path of the request.
func ParseRequestPath(r *http.Request) string {
	return r.URL.Path
}

// DataResponse generates body with data.
func DataResponse(resData ...any) gin.H {
	l := len(resData)
	if l > 0 {
		if l == 1 {
			return gin.H{
				"data": resData[0],
			}
		}

		return gin.H{
			"data": resData,
		}
	}

	return gin.H{
		"data": "",
	}
}
