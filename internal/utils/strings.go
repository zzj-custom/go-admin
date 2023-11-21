package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

// ToCamelCase 下换线转为驼峰
func ToCamelCase(s string) string {
	s = strings.ToLower(s)
	slice := strings.Split(s, "_")
	for i := 1; i < len(slice); i++ {
		slice[i] = cases.Title(language.Und).String(slice[i])
	}
	return strings.Join(slice, "")
}

// ToUnderscore 驼峰转为下划线
func ToUnderscore(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i != 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
