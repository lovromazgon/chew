package chew

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// spf13 / cobra
// Masterminds / sprig

const (
	TEMPLATE_SUFFIX = ".tmpl"
)

type Template struct {
	*template.Template
}

func New(name string) *Template {
	ct := &Template{template.New(name)}

	funcMap := template.FuncMap{
		"indentTemplate":  ct.IndentTemplate,
		"indentTemplates": ct.IndentTemplates,
		"repeat":          strings.Repeat,
	}

	ct.Funcs(funcMap)
	ct.Option("missingkey=error")
	return ct
}

func (ct *Template) ParseFolder(folderPath string) (*Template, error) {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, TEMPLATE_SUFFIX) {
			_, err = ct.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return err
	})

	return ct, err
}

// allows user to indent the template which is inserted
func (ct *Template) IndentTemplate(template string, data interface{}, parent interface{}, indentSize int) string {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		panic("nested template is not a map!")
	}
	dataMap["parent"] = parent

	buffer := new(bytes.Buffer)
	tmpl := ct.Lookup(template + TEMPLATE_SUFFIX)
	if tmpl == nil {
		panic(fmt.Sprintf("Could not find template '%s%s'", template, TEMPLATE_SUFFIX))
	}
	tmpl.Execute(buffer, data)

	return Indent(indentSize, buffer.String())
}

func (ct *Template) IndentTemplates(nestedTemplates interface{}, templateField string, parent interface{}, indentSize int) string {
	var nestedTemplatesSlice []interface{}
	nestedTemplatesSlice, ok := nestedTemplates.([]interface{})
	if !ok {
		panic("nested templates are not an array!")
	}

	buffer := new(bytes.Buffer)
	for i, data := range nestedTemplatesSlice {
		var dataMap map[string]interface{}
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			panic("nested template is not a map!")
		}

		var tmpl string
		if templateValue, ok := dataMap[templateField]; ok {
			if templateString, ok := templateValue.(string); ok {
				tmpl = templateString
			} else {
				panic(fmt.Sprintf("Field %s in %+v is not of type string!", templateField, data))
			}
		} else {
			panic(fmt.Sprintf("Could not find field %s in %+v", templateField, data))
		}

		buffer.WriteString(ct.IndentTemplate(tmpl, data, parent, indentSize))
		if i < len(nestedTemplatesSlice)-1 {
			buffer.WriteString("\n")
			buffer.WriteString(strings.Repeat(" ", indentSize))
		}
	}
	return buffer.String()
}

func Indent(identSize int, a string) string {
	scanner := bufio.NewScanner(strings.NewReader(a))
	var buffer bytes.Buffer

	i := 0
	for scanner.Scan() {
		if i > 0 {
			buffer.WriteString(strings.Repeat(" ", identSize))
		}
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
		i++
	}
	return buffer.String()
}
