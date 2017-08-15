package chew

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"bitbucket.org/lovromazgon/chew/funcmap"
)

const (
	TEMPLATE_SUFFIX = ".tmpl"
)

type Template struct {
	*template.Template
	Functions funcmap.Functions

	injectFuncsOnce sync.Once
}

func New(name string) *Template {
	ct := &Template{
		Template:  template.New(name),
		Functions: funcmap.Global,
	}
	ct.InjectFunctions()
	ct.Option("missingkey=error")
	return ct
}

func (ct *Template) InjectFunctions() {
	ct.injectFuncsOnce.Do(func() {
		ct.Functions = ct.Functions.MustAddFunc(funcmap.NewFunc(
			func() map[string]interface{} {
				return map[string]interface{}{
					"version":        Version,
					"version_date":   VersionDate,
					"execution_date": time.Now().Format("02.01.2006"),
					"execution_time": time.Now().Format("15:04"),
				}
			},
			"chew",
			"TODO",
			"TODO",
		)).MustAddFunc(funcmap.NewFunc(
			ct.IndentTemplate,
			"indentTemplate",
			"Use indentTemplate to execute a child template and indent the content of the template with spaces."+
				" This function takes 3 parameters:\n"+
				"- template string : the name of the nested template\n"+
				"- data interface{} : the data on which the nested template will be evaluated\n"+
				"- parent interface{} : the parent who calls the nested template (will be available in the nested template in field .parent)\n"+
				"- indentSize int : number of spaces to indent this template",
			"{{ indentTemplate .child.template_field .child . 2 }}",
		)).MustAddFunc(funcmap.NewFunc(
			ct.IndentTemplates,
			"indentTemplates",
			"TODO",
			"TODO",
		)).MustAddFunc(funcmap.NewFunc(
			ct.Plugins,
			"plugins",
			"TODO",
			"TODO",
		))

		ct.Funcs(ct.Functions.FuncMap())
	})
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

func (ct *Template) ExecuteChewable(w Writer, c Chewable) error {
	for _, cd := range c.Data {
		for tmpl, out := range cd.Templates {
			w.SetOut(out)
			err := ct.Template.ExecuteTemplate(w, tmpl+TEMPLATE_SUFFIX, prepareData(c, cd))
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

	for k, v := range c.Global {
		data[k] = v
	}
	for k, v := range cd.Local {
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

	return funcmap.Indent(indentSize, buffer.String())
}

func (ct *Template) IndentTemplates(nestedTemplates interface{}, templateField string, parent interface{}, indentSize int) string {
	return ct.Plugins(nestedTemplates, "", templateField, parent, indentSize)
}

func (ct *Template) Plugins(pluginsRaw interface{}, insertPoint, templateField string, parent interface{}, indentSize int) string {
	if pluginsRaw == nil {
		// it can be tha the key doesn't exist
		return ""
	}

	var pluginsSlice []interface{}
	pluginsSlice, ok := pluginsRaw.([]interface{})
	if !ok {
		panic("nested templates are not an array!")
	}

	buffer := new(bytes.Buffer)
	for _, data := range pluginsSlice {
		var dataMap map[string]interface{}
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			panic("nested template is not a map!")
		}

		var tmpl string
		if templateInterface, ok := dataMap[templateField]; ok {
			if templateString, ok := templateInterface.(string); ok {
				tmpl = templateString
			} else if templateSlice, ok := templateInterface.(map[string]interface{}); ok {
				if templateRaw, ok := templateSlice[insertPoint]; ok {
					if templateString, ok := templateRaw.(string); ok {
						tmpl = templateString
					} else {
						panic(fmt.Sprintf("Field %s.%s is not a string or slice!", templateField, insertPoint))
					}
				} else {
					// point of insert not found for this plugin - no problem
					continue
				}
			} else {
				panic(fmt.Sprintf("Field %s is not a string or slice!", templateField))
			}
		} else if insertPoint == "" {
			// no insert point, this means it is not a plugin but a nested template
			panic(fmt.Sprintf("Could not find field %s", templateField))
		} else {
			// we are searching for an insert point - not needed to be found
			continue
		}

		buffer.WriteString(ct.IndentTemplate(tmpl, data, parent, indentSize))
		buffer.WriteString("\n")
	}
	return strings.TrimRight(buffer.String(), "\n")
}
