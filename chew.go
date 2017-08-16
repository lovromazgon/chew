// Chew is a lean and fast CLI wrapper for Go text/template with additional functions.
// The template generation in Chew is data-centric, meaning that the input data dictates
// which templates will be generated and with which data.
//
// You can define a folder where you store all of your templates and Chew will parse
// all files in the folder and sub-folders when you run it. When running chew you also
// define the input data in JSON format, which will be used to generate the output.
//
// For example if we have the folder /templates:
//   ▾ templates
//     ▾ special
//         bearing.tmpl
//       case.tmpl
//       fidget_spinner.tmpl
//
// Content of fidget_spinner.tmpl:
//   Spinner model: {{ .model }}
//   Spinner year of construction: {{ .construction.year }}
//
//   Parts:
//   {{ indentTemplate "case" .case . 2 }}
//
// Content of case.tmpl:
//   Case type: {{ .parent.model }}.{{ .type }}
//   Main Bearings:
//   {{ plugins .bearings "template" "main" . 2 }}
//
//   Outer Bearings:
//   {{ plugins .bearings "template" "outer" . 2 }}
//
//
// Content of bearing.tmpl:
//
package chew

import (
	"time"

	_ "bitbucket.org/lovromazgon/chew-plsql"
	"bitbucket.org/lovromazgon/chew/funcmap"
)

var (
	Version     string
	VersionDate string
)

func init() {
	funcmap.AddFunc(&funcmap.Func{
		Func: func() map[string]interface{} {
			return map[string]interface{}{
				"version":        Version,
				"version_date":   VersionDate,
				"execution_date": time.Now().Format("02.01.2006"),
				"execution_time": time.Now().Format("15:04"),
			}
		},
		Doc: funcmap.FuncDoc{
			Name:    "chew",
			Text:    "Makes some general information about Chew available in templates",
			NestedFuncs: []funcmap.FuncDoc {
				{
					Name: "version",
					Text: "Returns the current version of Chew",
					Example: "{{ chew.version }}",
				},
				{
					Name: "version_date",
					Text: "Returns the date of the release of Chew",
					Example: "{{ chew.version_date }}",
				},
				{
					Name: "execution_date",
					Text: "Returns the date of execution (today)",
					Example: "{{ chew.execution_date }}",
				},
				{
					Name: "execution_time",
					Text: "Returns the time of execution (now)",
					Example: "{{ chew.execution_time }}",
				},
			},
		},
	})
}
