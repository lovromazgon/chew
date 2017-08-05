package chew

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/lovromazgon/chew/plsql"
)

func NewFuncMap(ct *Template) map[string]interface{} {
	funcMap := map[string]interface{}{
		"chew": func() map[string]interface{} {
			return map[string]interface{}{
				"version":        VERSION,
				"version_date":   VERSION_DATE,
				"execution_date": time.Now().Format("02.01.2006"),
				"execution_time": time.Now().Format("15:04"),
			}
		},
		"indentTemplate":  ct.IndentTemplate,
		"indentTemplates": ct.IndentTemplates,
		"plugins":         ct.Plugins,
		"maxLength":       MaxLength,
		"offset":          Offset,
		"exists":          Exists,
	}

	for k, v := range plsql.NewFuncMap() {
		if _, ok := funcMap[k]; ok {
			panic(fmt.Sprintf("Global function map already contains function %s", k))
		}
		funcMap[k] = v
	}

	return funcMap
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
