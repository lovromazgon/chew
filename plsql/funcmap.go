package plsql

import (
	"fmt"
	"strings"
	"sync"
)

var (
	plsqlf     *PlsqlFuncMap
	plsqlfInit sync.Once
)

func PlsqlNS() *PlsqlFuncMap {
	plsqlfInit.Do(func() { plsqlf = &PlsqlFuncMap{} })
	return plsqlf
}

func NewFuncMap() map[string]interface{} {
	funcMap := map[string]interface{}{
		"plsql": PlsqlNS,
	}

	return funcMap
}

type PlsqlFuncMap struct {}

func (*PlsqlFuncMap) ParameterType(paramType string) string {
	switch strings.ToUpper(paramType) {
	case "IN":
		return "IN    "
	case "OUT":
		return "   OUT"
	case "IN OUT":
		return "IN OUT"
	default:
		panic(fmt.Sprintf("Unknown parameter type: %s", paramType))
	}
}
