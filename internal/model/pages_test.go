package model

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageDynamicURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/path", nil)

	var dynamic = PageDynamicURL{}
	dynamic.Parse(req)

	assert.Equal(t, "example.com", dynamic.Convert("{host}"))
	assert.Equal(t, "example.com", dynamic.Convert("{hostname}"))
	assert.Equal(t, "http://example.com/path", dynamic.Convert("{href}"))
	assert.Equal(t, "http://example.com", dynamic.Convert("{origin}"))
	assert.Equal(t, "/path", dynamic.Convert("{pathname}"))
	assert.Equal(t, "80", dynamic.Convert("{port}"))
	assert.Equal(t, "http:", dynamic.Convert("{protocol}"))
	assert.Equal(t, "{hash}", dynamic.Convert("{hash}"))
	assert.Equal(t, "{search}", dynamic.Convert("{search}"))
	assert.Equal(t, "https://example.com:8080/redirect=/path", dynamic.Convert("https://{hostname}:8080/redirect={pathname}"))
}
