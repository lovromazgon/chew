package chew

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"bitbucket.org/lovromazgon/chew/funcmap"
)

var (
	templateSuffix = ".tmpl"
)

// Template wraps *text/template.Template and adds some additional functionality. It should be
// created via chew.New. If a manual instantiation is used the method InjectFunctions should be
// called before processing templates to be able to use all custom chew functions.
type Template struct {
	*template.Template
	Functions funcmap.Functions

	injectFuncsOnce sync.Once
}

// New creates a new chew.Template with the provided name and injects the template functions.
func New(name string) *Template {
	ct := &Template{
		Template:  template.New(name),
		Functions: funcmap.Global,
	}
	ct.InjectFunctions()
	ct.Option("missingkey=error")
	return ct
}

// InjectFunctions adds template specific functions to Template.Functions and adds the
// function map to the template.
func (ct *Template) InjectFunctions() {
	ct.injectFuncsOnce.Do(func() {
		ct.Functions = ct.Functions.MustAddFunc(&funcmap.Func{
			Func: ct.IndentTemplate,
			Doc: funcmap.FuncDoc{
				Name: "indentTemplate",
				Text: "Use indentTemplate to execute a child template and indent the content of the template with spaces." +
					" This function takes 4 parameters:\n" +
					"- template string    : the name of the nested template\n" +
					"- data interface{}   : the data on which the nested template will be evaluated\n" +
					"- parent interface{} : the parent who calls the nested template (will be available in the nested template in field .parent)\n" +
					"- indentSize int     : number of spaces to indent this template",
				Example: "{{ indentTemplate .child.template_field .child . 2 }}",
			},
		}).MustAddFunc(&funcmap.Func{
			Func: ct.IndentTemplates,
			Doc: funcmap.FuncDoc{
				Name: "indentTemplates",
				Text: "Use indentTemplate to execute multiple child templates and indent the content with spaces." +
					"This function takes 4 parameters:\n" +
					"- nestedTemplates []interface{} : the array whic contains the nested templates\n" +
					"- templateField string          : the name of the field which contains the name of the template to be used\n" +
					"- parent interface{}            : the parent who calls the nested template (will be available in the nested template in field .parent)\n" +
					"- indentSize int                : number of spaces to indent this template",
				Example: "{{ indentTemplates .childArray \"template_field\" . 2 }}",
			},
		}).MustAddFunc(&funcmap.Func{
			Func: ct.Plugins,
			Doc: funcmap.FuncDoc{
				Name: "plugins",
				Text: "Use plugins to evaluate 0 or more plugins which can choose which template to call on which insert point." +
					"This function takes 5 parameters:\n" +
					"- plugins []interface{} : the array which contains the plugins\n" +
					"- insertPoint string    : the insertion point defined in the parent template where plugins can insert some content\n" +
					"- templateField string  : the name of the field which contains the name of the template to be used\n" +
					"- parent interface{}    : the parent who calls the nested template (will be available in the nested template in field .parent)\n" +
					"- indentSize int        : number of spaces to indent this template",
				Example: "{{ plugins .pluginsArray \"insert_point_1\" \"template_field\" . 2 }}",
			},
		})

		ct.Funcs(ct.Functions.FuncMap())
	})
}

// ParseFolder recursively walks through the provided folder path and parses every template
// it can find with the template suffix.
func (ct *Template) ParseFolder(folderPath string) (*Template, error) {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, templateSuffix) {
			_, err = ct.ParseFiles(path)
		}
		return err
	})

	return ct, err
}

// ExecuteChewable loops through all ChewableData in Chewable and executes every Template defined
// in ChewableData.Templates. The output is written to the supplied chew.Writer, which also gets notified
// about the desired output filename before every template execution.
// If template.Template.ExecuteTemplate returns an error the execution stops and returns it.
func (ct *Template) ExecuteChewable(w Writer, c Chewable) error {
	for _, cd := range c.Data {
		for tmpl, out := range cd.Templates {
			w.SetOut(out)
			err := ct.Template.ExecuteTemplate(w, tmpl+templateSuffix, prepareData(c, cd))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Merges local data from ChewableData and global data from Chewable.
// If a key exists in both global and local, then local is used.
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

// IndentTemplate is similar to the built-in template function with the additional functionality
// of indenting the content of the template and accessing parent data.
// It takes the name of the template, the data which will be sent to the template when executing it,
// the parent data which will be accessible in the executed template and the number of spaces which will
// indent the content of the processed template.
func (ct *Template) IndentTemplate(template string, data interface{}, parent interface{}, indentSize int) string {
	dataMap, err := ToMap(data)
	if err != nil {
		panic(err)
	}

	dataMap["parent"] = parent

	buffer := new(bytes.Buffer)
	tmpl := ct.Lookup(template + templateSuffix)
	if tmpl == nil {
		panic(fmt.Sprintf("Could not find template '%s%s'", template, templateSuffix))
	}
	if err := tmpl.Execute(buffer, dataMap); err != nil {
		panic(err)
	}

	return Indent(indentSize, buffer.String())
}

// IndentTemplates is similar to IndentTemplate, only that it processes all templates defined in a slice.
// It takes the slice of objects which carry the data about the nested templates, the field in the nested template
// object which carries the name of the template, the parent data which will be accessible in the executed template
// and the number of spaces which will indent the content of the processed template.
func (ct *Template) IndentTemplates(nestedTemplates interface{}, templateField string, parent interface{}, indentSize int) string {
	return ct.Plugins(nestedTemplates, "", templateField, parent, indentSize)
}

// Plugins is similar to IndentTemplates, only that it skips nested templates which don't define a template
// for the insertion point. It takes the slice of objects which carry the data about the plugins, the insertion
// point where the plugin will be inserted, the field in the plugin object which carries the name of the template,
// the parent data which will be accessible in the executed template and the number of spaces which will indent
// the content of the processed template.
func (ct *Template) Plugins(pluginsRaw interface{}, insertPoint, templateField string, parent interface{}, indentSize int) string {
	if pluginsRaw == nil {
		// it can be tha the key doesn't exist
		return ""
	}

	var pluginsSlice []interface{}
	pluginsSlice, ok := pluginsRaw.([]interface{})
	if !ok {
		panic("nested templates are not a slice!")
	}

	buffer := new(bytes.Buffer)
	for _, data := range pluginsSlice {
		dataMap, err := ToMap(data)
		if err != nil {
			panic(err)
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
