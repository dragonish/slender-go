package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDomain(t *testing.T) {
	assert.Equal(t, "example.com", GetDomain("test.example.com"))
	assert.Equal(t, "test.example.com", GetDomain("abc.test.example.com"))
	assert.Equal(t, "example.com", GetDomain("example.com"))
	assert.Equal(t, "192.168.1.1", GetDomain("192.168.1.1"))
	assert.Equal(t, "127.0.0.1", GetDomain("127.0.0.1"))
	assert.Equal(t, "localhost", GetDomain("localhost"))
	assert.Equal(t, "::1", GetDomain("::1"))
}
