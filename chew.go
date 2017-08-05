package chew

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"io"
)

// spf13 / cobra
// Masterminds / sprig

const (
	TEMPLATE_SUFFIX = ".tmpl"
	VERSION = "0.0.1"
	VERSION_DATE = "05.08.2017"
)

type Template struct {
	*template.Template
}

func New(name string) *Template {
	ct := &Template{template.New(name)}
	ct.Funcs(NewFuncMap(ct)).Option("missingkey=error")
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

func (ct *Template) ExecuteChewable(w io.Writer, c Chewable) error {
	for _,cd := range c.Data {
		for _,tmpl := range cd.Templates {
			fmt.Println("---------------------------")
			err := ct.Template.ExecuteTemplate(w, tmpl + TEMPLATE_SUFFIX, prepareData(c, cd))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// merges local data from ChewableData and global data from Chewable
// if a key exists in both global and local, then local is used
func prepareData(c Chewable, cd ChewableData) map[string]interface{} {
	data := make(map[string]interface{})

	for k,v := range c.Global {
		data[k] = v
	}
	for k,v := range cd.Local {
		data[k] = v
	}

	return data
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
	if err := tmpl.Execute(buffer, data); err != nil {
		panic(err)
	}

	return Indent(indentSize, buffer.String())
}

func (ct *Template) IndentTemplates(nestedTemplates interface{}, templateField string, parent interface{}, indentSize int) string {
	if nestedTemplates == nil {
		// it can be tha the key doesn't exist
		return ""
	}

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
		}
	}
	return buffer.String()
}