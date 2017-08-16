package chew

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"log"
)

func ExampleTemplate_ExecuteChewable() {
	dataRaw, _ := ioutil.ReadFile("test/data/example.json")

	chewable := &Chewable{}
	err := json.Unmarshal(dataRaw, chewable)
	if err != nil {
		log.Fatal(err)
	}

	template := New("main")
	template.ParseFolder("test/templates")
	template.ExecuteChewable(WriterWrapper{os.Stdout}, *chewable)
}

func TestTemplate_IndentTemplate(t *testing.T) {
	data := map[string]interface{}{
		"local_var": "test",
	}
	parent := map[string]interface{}{
		"local_var": "parent_test",
	}

	template := New("main")
	template.ParseFolder("test/templates")
	actual := template.IndentTemplate("test_indentTemplate", data, parent, 4)
	expected := `    my local variable:'test'
    local variable from my parent:'parent_test'`

	assert.Equal(t, expected, actual)
}

func TestTemplate_IndentTemplates(t *testing.T) {
	data := map[string]interface{}{
		"nested": []interface{}{
			map[string]interface{}{
				"name":      "First",
				"template":  "test_plugins_plugin1_ger",
				"template2": "test_plugins_plugin1_ita",
			},
			map[string]interface{}{
				"name":     "Second",
				"template": "test_plugins_plugin2",
			},
		},
	}

	template := New("main")
	template.ParseFolder("test/templates")

	// for insertion point 1 only the first plugin should be found
	actual := template.IndentTemplates(data["nested"], "template", data, 2)
	expected := `  Plugin Nummer eins:
  I got inserted by 'First'
  Plugin Nummer zwei:
  I got inserted by 'Second'`
	assert.Equal(t, expected, actual)

	assert.Panics(t, func() {
		// there is no templat2 field in the second template - panic
		template.IndentTemplates(data["nested"], "template2", data, 2)
	})
}

func TestTemplate_Plugins(t *testing.T) {
	data := map[string]interface{}{
		"plugins": []interface{}{
			map[string]interface{}{
				"name": "First",
				"template": map[string]interface{}{
					"insertion_point_1": "test_plugins_plugin1_ita",
					"insertion_point_2": "test_plugins_plugin1_ger",
				},
			},
			map[string]interface{}{
				"name": "Second",
				"template": map[string]interface{}{
					"insertion_point_2": "test_plugins_plugin2",
				},
			},
		},
	}

	template := New("main")
	template.ParseFolder("test/templates")

	// for insertion point 1 only the first plugin should be found
	actualIns1 := template.Plugins(data["plugins"], "insertion_point_1", "template", data, 1)
	expectedIns1 := ` Plugin numero uno:
 I got inserted by 'First'`
	assert.Equal(t, expectedIns1, actualIns1)

	// for insertion point 2 both plugins should be found
	actualIns2 := template.Plugins(data["plugins"], "insertion_point_2", "template", data, 3)
	expectedIns2 := `   Plugin Nummer eins:
   I got inserted by 'First'
   Plugin Nummer zwei:
   I got inserted by 'Second'`
	assert.Equal(t, expectedIns2, actualIns2)

	// for insertion point 3 no plugins should be found
	actualIns3 := template.Plugins(data["plugins"], "insertion_point_3", "template", data, 5)
	expectedIns3 := ""

	assert.Equal(t, expectedIns3, actualIns3)
}
