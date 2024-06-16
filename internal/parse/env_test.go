package parse

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnvVars(t *testing.T) {
	os.Setenv("SLENDER_ACCESS_PWD", "abc456")
	os.Setenv("SLENDER_ADMIN_PWD", "123efg")
	os.Setenv("SLENDER_LOG_LEVEL", "Error")
	os.Setenv("SLENDER_PORT", "10086")

	envs := parseEnvVars()
	assert.Equal(t, "abc456", envs.AccessPassword)
	assert.Equal(t, "123efg", envs.AdminPassword)
	assert.Equal(t, "Error", envs.LogLevel)
	assert.Equal(t, uint16(10086), envs.Port)
}
