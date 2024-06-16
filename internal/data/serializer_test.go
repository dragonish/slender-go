package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelCaseToSnakeCase(t *testing.T) {
	assert.Equal(t, "hello_world", camelCaseToSnakeCase("HelloWorld"))
	assert.Equal(t, "hello_world", camelCaseToSnakeCase("helloWorld"))
	assert.Equal(t, "hello_world", camelCaseToSnakeCase("hello_world"))
}

func TestStructToMap(t *testing.T) {
	type testStruct struct {
		ID      int64
		Name    string `json:"na"`
		Note    string `json:"not" db:"no"`
		ListA   []string
		ListB   []string
		HeightA int64 `json:"heightA,omitempty"`
		HeightB int64
		PointA  *int64
		PointB  *string `json:"pointB,omitempty"`
	}

	var num int64 = 3
	var val = testStruct{
		ID:      1,
		Name:    "user",
		ListA:   []string{"abc"},
		HeightB: 2,
		PointA:  &num,
	}

	m := StructToMap(val, "HeightB")

	assert.Contains(t, m, "id")
	assert.Contains(t, m, "na")
	assert.Contains(t, m, "no")
	assert.Contains(t, m, "list_a")
	assert.Contains(t, m, "list_b")
	assert.Contains(t, m, "height_a")
	assert.NotContains(t, m, "height_b")
	assert.Contains(t, m, "point_a")
	assert.NotContains(t, m, "point_b")
}

func TestGetSizeFromStr(t *testing.T) {
	assert.Equal(t, 30, GetSizeFromStr("30", 25, 10, 100))
	assert.Equal(t, 25, GetSizeFromStr("abc", 25, 10, 100))
	assert.Equal(t, 25, GetSizeFromStr("", 25, 10, 100))
	assert.Equal(t, 10, GetSizeFromStr("0", 25, 10, 100))
	assert.Equal(t, 10, GetSizeFromStr("-1", 25, 10, 100))
	assert.Equal(t, 100, GetSizeFromStr("1000", 25, 10, 100))
}

func TestGetPageFromStr(t *testing.T) {
	assert.Equal(t, 5, GetPageFromStr("5"))
	assert.Equal(t, 1, GetPageFromStr("0"))
	assert.Equal(t, 1, GetPageFromStr("-1"))
	assert.Equal(t, 1, GetPageFromStr("abc"))
	assert.Equal(t, 1, GetPageFromStr(""))
}

func TestIsRouteTruthy(t *testing.T) {
	assert.True(t, IsRouteTruthy("1"))
	assert.True(t, IsRouteTruthy("yes"))
	assert.True(t, IsRouteTruthy("true"))
	assert.True(t, IsRouteTruthy("on"))
	assert.False(t, IsRouteTruthy("0"))
	assert.False(t, IsRouteTruthy("no"))
	assert.False(t, IsRouteTruthy("false"))
	assert.False(t, IsRouteTruthy("off"))
	assert.False(t, IsRouteTruthy(""))
}

func TestIsRouteFalsy(t *testing.T) {
	assert.True(t, IsRouteFalsy("0"))
	assert.True(t, IsRouteFalsy("no"))
	assert.True(t, IsRouteFalsy("false"))
	assert.True(t, IsRouteFalsy("off"))
	assert.False(t, IsRouteFalsy("1"))
	assert.False(t, IsRouteFalsy("yes"))
	assert.False(t, IsRouteFalsy("true"))
	assert.False(t, IsRouteFalsy("on"))
	assert.False(t, IsRouteFalsy(""))
}

func TestDefference(t *testing.T) {
	slice1 := []int{2, 3, 4, 5}
	slice2 := []int{1, 3, 5, 7}

	res1 := Defference(slice1, slice2)
	res2 := Defference(slice2, slice1)

	assert.Len(t, res1, 2)
	assert.Contains(t, res1, 2)
	assert.Contains(t, res1, 4)
	assert.NotContains(t, res1, 3)
	assert.NotContains(t, res1, 5)
	assert.Len(t, res2, 2)
	assert.Contains(t, res2, 1)
	assert.Contains(t, res2, 7)
	assert.NotContains(t, res2, 3)
	assert.NotContains(t, res2, 5)
}

func TestInt16ToStringWithSign(t *testing.T) {
	assert.Equal(t, "+1", Int16ToStringWithSign(1))
	assert.Equal(t, "-1", Int16ToStringWithSign(-1))
	assert.Equal(t, "+0", Int16ToStringWithSign(0))
}
