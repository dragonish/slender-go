package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyString(t *testing.T) {
	var str = MyString(`example_keyword%\`)

	assert.Equal(t, `example\_keyword\%\\`, str.EscapeSpecialCharacters())
	assert.Equal(t, `%example\_keyword\%\\%`, str.LikeMatchingString())
}
