package chew

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"bitbucket.org/lovromazgon/chew/funcmap"
)

func init() {
	funcmap.AddFunc(&funcmap.Func{
		Func: Indent,
		Doc: funcmap.FuncDoc{
			Name:    "indent",
			Text:    "Indent prepends spaces to every line in a multi-line string",
			Example: "{{ indent 2 \"my beautiful\\n multiline string\" }}",
		},
	})
	funcmap.AddFunc(&funcmap.Func{
		Func: MaxLength,
		Doc: funcmap.FuncDoc{
			Name:    "maxLength",
			Text:    "MaxLength returns the length of the longest field in a []string",
			Example: "{{ maxLength .myStringSlice }}",
		},
	})
	funcmap.AddFunc(&funcmap.Func{
		Func: Offset,
		Doc: funcmap.FuncDoc{
			Name:    "offset",
			Text:    "Offset returns a string of blank spaces so that the input string reaches the input length",
			Example: "{{ offset 25 \"Need 4 spaces till 25\" }}",
		},
	})
	funcmap.AddFunc(&funcmap.Func{
		Func: Exists,
		Doc: funcmap.FuncDoc{
			Name:    "exists",
			Text:    "Exists returns true if field exists in a map, else false",
			Example: "{{ exists . \"my_field\" }}",
		},
	})
}

// Indent prepends spaces to every line in a multi-line string.
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

// Offset returns a string of blank spaces so that the input string reaches the input length.
func Offset(length int, a interface{}) string {
	aStr := fmt.Sprint(a)
	return strings.Repeat(" ", length-len(aStr))
}

// MaxLength returns the length of the longest field in a []string.
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

// Exists returns true if field exists in a map, else false.
func Exists(data map[string]interface{}, key string) bool {
	_, exists := data[key]
	return exists
}
