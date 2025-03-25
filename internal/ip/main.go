package ip

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetRealIP returns the real IP address of the client making the request.
//
// It first checks for X-Real-IP header, then X-Forwarded-For header, and
// finally falls back to ClientIP if both are empty.
func GetRealIP(ctx *gin.Context) (realIp string) {
	xri := ctx.GetHeader("X-Real-IP")
	if xri != "" {
		realIp = xri
		return
	}

	xff := ctx.GetHeader("X-Forwarded-For")
	ips := strings.Split(xff, ",")
	realIp = strings.TrimSpace(ips[len(ips)-1])

	if realIp == "" {
		realIp = ctx.ClientIP()
	}
	return
}
