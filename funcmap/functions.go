package funcmap

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func init() {
	AddFunc(NewFunc(
		Indent,
		"indent",
		"TODO",
		"TODO",
	))
	AddFunc(NewFunc(
		MaxLength,
		"maxLength",
		"TODO",
		"TODO",
	))
	AddFunc(NewFunc(
		Offset,
		"offset",
		"TODO",
		"TODO",
	))
	AddFunc(NewFunc(
		Exists,
		"exists",
		"TODO",
		"TODO",
	))
}

func Indent(identSize int, a string) string {
	if a == "" {
		return strings.Repeat(" ", identSize)
	}

	scanner := bufio.NewScanner(strings.NewReader(a))
	var buffer bytes.Buffer

	i := 0
	for scanner.Scan() {
		buffer.WriteString(strings.Repeat(" ", identSize))
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
		i++
	}

	return strings.TrimRight(buffer.String(), "\n")
}

func Offset(length int, a interface{}) string {
	aStr := fmt.Sprint(a)
	return strings.Repeat(" ", length-len(aStr))
}

func MaxLength(data interface{}) int {
	dataSlice, ok := data.([]string)
	if !ok {
		panic(fmt.Sprintf("Can't convert %+v to []string", data))
	}
	length := 0
	for _, str := range dataSlice {
		if len(str) > length {
			length = len(str)
		}
	}
	return length
}

func Exists(data map[string]interface{}, key string) bool {
	_, exists := data[key]
	return exists
}
