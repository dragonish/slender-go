package parse

import (
	"os"
	"path"
	"slender/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCLI(t *testing.T) {
	workDir, _ := os.Getwd()

	var base model.Flags = model.Flags{
		AccessPassword:  "123456",
		AdminPassword:   "654321",
		LogLevel:        "Info",
		Port:            8080,
		PerformanceMode: true,
		TokenAge:        30,
	}

	originalBase := base
	originalArgs := os.Args
	restore := func() {
		base = originalBase
		os.Args = originalArgs
	}
	defer restore()

	t.Run("Start command", func(t *testing.T) {
		defer restore()

		os.Args = []string{"slender", "--access_pwd", "", "--admin_pwd", "abcdef", "--port", "10086", "--performance=false"}
		res := parseCLI(&base)

		assert.Empty(t, res.ServiceConfig)
		assert.False(t, res.PerformanceMode)
		assert.Empty(t, res.AccessPassword)
		assert.Equal(t, "abcdef", res.AdminPassword)
		assert.Equal(t, uint16(30), res.TokenAge)
		assert.Equal(t, "Info", res.LogLevel)
		assert.Equal(t, uint16(10086), res.Port)
	})

	t.Run("Config file", func(t *testing.T) {
		defer restore()

		yamlPath := path.Join(workDir, "test.yaml")
		os.Remove(yamlPath)

		f, _ := os.Create(yamlPath)
		defer os.Remove(yamlPath)
		defer f.Close()
		f.Write([]byte("access_password: ''\nadmin_password: abcdef\nport: 10086\n"))

		base.ServiceConfig = "./test.yaml"
		os.Args = []string{"slender"}
		res := parseCLI(&base)

		assert.Equal(t, "./test.yaml", res.ServiceConfig)
		assert.True(t, res.PerformanceMode)
		assert.Empty(t, res.AccessPassword)
		assert.Equal(t, "abcdef", res.AdminPassword)
		assert.Equal(t, uint16(30), res.TokenAge)
		assert.Equal(t, "Info", res.LogLevel)
		assert.Equal(t, uint16(10086), res.Port)
	})

	t.Run("Start command & Config file", func(t *testing.T) {
		defer restore()

		yamlPath := path.Join(workDir, "test.yaml")
		os.Remove(yamlPath)

		f, _ := os.Create(yamlPath)
		defer os.Remove(yamlPath)
		defer f.Close()
		f.Write([]byte("access_password: ''\nadmin_password: abcdef\nport: 10086\n"))

		os.Args = []string{"slender", "--access_pwd", "112233", "--port", "10010", "--config", "./test.yaml"}
		res := parseCLI(&base)

		assert.Equal(t, "./test.yaml", res.ServiceConfig)
		assert.True(t, res.PerformanceMode)
		assert.Equal(t, "112233", res.AccessPassword)
		assert.Equal(t, "abcdef", res.AdminPassword)
		assert.Equal(t, uint16(30), res.TokenAge)
		assert.Equal(t, "Info", res.LogLevel)
		assert.Equal(t, uint16(10010), res.Port)
	})
}
