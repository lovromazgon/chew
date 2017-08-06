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
// Custom functions:
//   {{-/* plugins [name of template] [data] [parent] [indent] */}}
//   {{-/* indentTemplate [name of template] [data] [parent] [indent] */}}
//
package chew

const (
	VERSION      = "0.0.1"
	VERSION_DATE = "05.08.2017"
)
