package chew

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"bitbucket.org/lovromazgon/plsql-templator"
)

func TestMain(m *testing.M) {
	dataRaw, err := ioutil.ReadFile("data/data.json")
	if err != nil {
		panic(err)
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(dataRaw, &data); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", data)

	myChew := New("main")
	_, err = myChew.ParseFolder("plsql-templates/common")

	if err != nil {
		panic(err)
	}

	fmt.Println("---------------------------")
	err = myChew.ExecuteTemplate(os.Stdout, "package_spec"+templator.TEMPLATE_SUFFIX, data["data"])
	if err != nil {
		panic(err)
	}

	/*
		json.Unmarshal()
		data := NewPackage("my_package", "my_schema", "This is a package for testing", "package")

		data.AddPlugin(PKG_BEFORE_PROCEDURES, NewBlockComment(fmt.Sprintf("PLSQL-Templator v0.0.1\nDate: %s", time.Now().String()), true))

		plugin := NewProcedure("my_plug", "My plugin procedure with parameters", "procedure")
		plugin.AddParameter(&Parameter{Name: "what", DataType: "VARCHAR2(10)", Type: IN})
		plugin.AddParameter(&Parameter{Name: "is", DataType: "NUMBER(10)", Type: IN})
		plugin.AddParameter(&Parameter{Name: "this", DataType: "TIMESTAMP", Type: IN, Default: "'this'"})

		data.AddPlugin(PKG_PROCEDURE, NewProcedure("plugged_in1", "Just a regular procedure", "procedure"))
		data.AddPlugin(PKG_BEFORE_PROCEDURES, NewProcedure("plugged_in2", "Same es plugged_in1", "procedure"))
		data.AddPlugin(PKG_PROCEDURE, plugin)

		data.AddPlugin(PKG_PROCEDURE, NewProcedure("proc_1", "Proc1 without params", "procedure"))
		proc2 := NewProcedure("proc_2", "This proc has some params", "procedure")
		data.AddPlugin(PKG_PROCEDURE, proc2)
		proc2.AddParameter(&Parameter{Name: "first_name", DataType: "VARCHAR2(10)", Type: IN})
		proc2.AddParameter(&Parameter{Name: "age", DataType: "NUMBER(8)", Type: INOUT})
		proc2.AddParameter(&Parameter{Name: "my_extremely_long_parameter_30", DataType: "TIMESTAMP", Type: INOUT, Default: "'YES'"})

		templates := templator.New("main")

		_, err := templates.ParseFolder("templates")

		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		data.Prepare(nil)
		data.SwitchModeRecursive(SPEC)
		err = templates.ExecuteTemplate(os.Stdout, data.Template()+templator.TEMPLATE_SUFFIX, data)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	*/
}
