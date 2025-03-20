package ip

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRealIP(ctx *gin.Context) (realIp string) {
	xff := ctx.GetHeader("X-Forwarded-For")
	ips := strings.Split(xff, ",")
	realIp = strings.TrimSpace(ips[len(ips)-1])

	if realIp == "" {
		realIp = ctx.ClientIP()
	}
	return
}
