package ip

import (
	"net/http"
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

// GetProtocol returns the protocol (http or https) used for the request.
func GetProtocol(r *http.Request) (proto string) {
	xfp := r.Header.Get("X-Forwarded-Proto")
	if xfp == "http" || xfp == "https" {
		proto = xfp
		return
	}

	proto = "http"
	if r.TLS != nil {
		proto = "https"
	}
	return
}
