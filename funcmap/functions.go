package funcmap

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func init() {
	AddFunc(&Func{
		Func: Indent,
		Doc: FuncDoc{
			Name:    "indent",
			Text:    "With indent you can prepend spaces to a multi-line string.",
			Example: "{{ indent 2 \"my beautiful\\n multiline string\" }}",
		},
	})
	AddFunc(&Func{
		Func: MaxLength,
		Doc: FuncDoc{
			Name:    "maxLength",
			Text:    "Returns the length of the longest field in a string slice",
			Example: "{{ maxLength .myStringSlice }}",
		},
	})
	AddFunc(&Func{
		Func: Offset,
		Doc: FuncDoc{
			Name:    "offset",
			Text:    "Returns the offset in blank spaces so that the input string reaches the input length",
			Example: "{{ offset 25 \"Need 4 spaces till 25\" }}",
		},
	})
	AddFunc(&Func{
		Func: Exists,
		Doc: FuncDoc{
			Name:    "exists",
			Text:    "Returns true if field exists in map, else false",
			Example: "{{ exists . \"my_field\" }}",
		},
	})
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
