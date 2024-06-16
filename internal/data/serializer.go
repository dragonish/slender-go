package data

import (
	"fmt"
	"reflect"
	"regexp"
	"slender/internal/logger"
	"slices"
	"strconv"
	"strings"
)

var matchFirstCap *regexp.Regexp
var matchAllCap *regexp.Regexp

func init() {
	matchFirstCap = compileReg("(.)([A-Z][a-z]+)")
	matchAllCap = compileReg("([a-z0-9])([A-Z])")
}

func compileReg(pattern string) *regexp.Regexp {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		logger.Fatal("error compiling regular expression", err, "pattern", pattern)
	}
	return reg
}

// camelCaseToSnakeCase converts the camelCase to the snakeCase.
func camelCaseToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)
	return snake
}

// StructToMap converts the struct to the map.
//
// Allows set filter list of field names.
// If the field names match, they will not be output to the result map.
//
// Map's key value be read from the db, json and attribute label,
// and auto convert to the snakeCase.
//
// If the field in the struct is a pointer type,
// the value of the reference is solved.
// If the pointer is nil, it will not be included in the map.
func StructToMap(input any, filterList ...string) map[string]interface{} {
	m := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typeOfInput := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typeOfInput.Field(i)
		if len(filterList) > 0 && slices.Contains(filterList, field.Name) {
			continue
		}

		value := val.Field(i)

		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			//? value.Elem() returns the value of the reference is solved.
			value = value.Elem()
		}

		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			dbParts := strings.Split(dbTag, ",")
			m[camelCaseToSnakeCase(dbParts[0])] = value.Interface()
		} else {
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" {
				tagParts := strings.Split(jsonTag, ",")
				m[camelCaseToSnakeCase(tagParts[0])] = value.Interface()
			} else {
				m[camelCaseToSnakeCase(field.Name)] = value.Interface()
			}
		}
	}

	return m
}

// GetSizeFromStr returns the number of size per page.
func GetSizeFromStr(source string, defaultValue, min, max int) int {
	limit, err := strconv.Atoi(source)
	if err != nil {
		limit = defaultValue
	}

	if limit < min {
		limit = min
	} else if limit > max {
		limit = max
	}

	return limit
}

// GetPageFromStr returns the number of pages.
//
// The minimum and default value are 1.
func GetPageFromStr(source string) int {
	page, err := strconv.Atoi(source)
	if err != nil {
		page = 1
	}

	if page < 1 {
		page = 1
	}

	return page
}

// IsRouteTruthy returns true when param is truthy.
func IsRouteTruthy(param string) bool {
	switch param {
	case "1", "yes", "true", "on":
		return true
	default:
		return false
	}
}

// IsRouteFalsy returns true when param is falsy.
func IsRouteFalsy(param string) bool {
	switch param {
	case "0", "no", "false", "off":
		return true
	default:
		return false
	}
}

// Defference returns left - right difference slice.
//
// The element that exists on the left but not on the right.
func Defference[T comparable](left, right []T) []T {
	diff := make([]T, 0)
	for i := 0; i < len(left); i++ {
		if !slices.Contains(right, left[i]) {
			diff = append(diff, left[i])
		}
	}
	return diff
}

// Int16ToStringWithSign converts an int16 to a string, adding a '+' character for positive or zero values.
func Int16ToStringWithSign(num int16) string {
	if num >= 0 {
		return fmt.Sprintf("+%d", num)
	}
	return strconv.Itoa(int(num))
}
