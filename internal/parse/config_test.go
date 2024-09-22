package parse

import (
	"os"
	"path"
	"slender/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceConfigHandler(t *testing.T) {
	var base model.Flags
	adminPassword := "abcdef"
	port := uint16(10086)
	var conf model.ServiceConfig = model.ServiceConfig{
		AdminPassword: &adminPassword,
		Port:          &port,
	}

	serviceConfigHandler(&base, &conf)

	assert.Empty(t, base.AccessPassword)
	assert.Empty(t, base.PerformanceMode)
	assert.Empty(t, base.TokenAge)
	assert.Empty(t, base.LogLevel)
	assert.Equal(t, port, base.Port)
	assert.Equal(t, adminPassword, base.AdminPassword)
}

func TestParseServiceConfig(t *testing.T) {
	t.Run("Content of YAML file", func(t *testing.T) {
		var content = []byte("access_password: ''\nadmin_password: abcdef\nport: 10086\n")

		conf, _ := parseServiceConfig(&content, false)
		assert.Empty(t, *conf.AccessPassword)
		assert.Equal(t, "abcdef", *conf.AdminPassword)
		assert.Nil(t, conf.LogLevel)
		assert.Equal(t, uint16(10086), *conf.Port)
		assert.Nil(t, conf.PerformanceMode)
		assert.Nil(t, conf.TokenAge)
	})

	t.Run("Content of JSON file", func(t *testing.T) {
		var content = []byte(`{
			"accessPassword": "",
			"adminPassword": "abcdef",
			"port": 10086
			}`)

		conf, _ := parseServiceConfig(&content, true)
		assert.Empty(t, *conf.AccessPassword)
		assert.Equal(t, "abcdef", *conf.AdminPassword)
		assert.Nil(t, conf.LogLevel)
		assert.Equal(t, uint16(10086), *conf.Port)
		assert.Nil(t, conf.PerformanceMode)
		assert.Nil(t, conf.TokenAge)
	})
}

func TestLoadServiceConfig(t *testing.T) {
	workDir, _ := os.Getwd()

	t.Run("YAML file", func(t *testing.T) {
		yamlPath := path.Join(workDir, "test.yaml")
		os.Remove(yamlPath)

		f, _ := os.Create(yamlPath)
		defer os.Remove(yamlPath)
		defer f.Close()
		f.Write([]byte("access_password: ''\nadmin_password: abcdef\nport: 10086\n"))

		var base model.Flags = model.Flags{
			AccessPassword:  "123456",
			AdminPassword:   "654321",
			LogLevel:        "Info",
			Port:            8080,
			PerformanceMode: true,
			TokenAge:        30,
		}
		loadServiceConfig(&base, yamlPath)
		assert.Empty(t, base.AccessPassword)
		assert.Equal(t, "abcdef", base.AdminPassword)
		assert.Equal(t, "Info", base.LogLevel)
		assert.Equal(t, uint16(10086), base.Port)
		assert.True(t, base.PerformanceMode)
		assert.Equal(t, uint16(30), base.TokenAge)
	})

	t.Run("JSON file", func(t *testing.T) {
		jsonPath := path.Join(workDir, "test.json")
		os.Remove(jsonPath)

		f, _ := os.Create(jsonPath)
		defer os.Remove(jsonPath)
		defer f.Close()
		f.Write([]byte(`{
			"accessPassword": "",
			"adminPassword": "abcdef",
			"port": 10086
			}`))

		var base model.Flags = model.Flags{
			AccessPassword:  "123456",
			AdminPassword:   "654321",
			LogLevel:        "Info",
			Port:            8080,
			PerformanceMode: true,
			TokenAge:        30,
		}
		loadServiceConfig(&base, jsonPath)
		assert.Empty(t, base.AccessPassword)
		assert.Equal(t, "abcdef", base.AdminPassword)
		assert.Equal(t, "Info", base.LogLevel)
		assert.Equal(t, uint16(10086), base.Port)
		assert.True(t, base.PerformanceMode)
		assert.Equal(t, uint16(30), base.TokenAge)
	})
}
