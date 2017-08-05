package chew

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChewable_UnmarshalJSON(t *testing.T) {
	dataRaw, err := ioutil.ReadFile("test/data/chewable_test.json")
	assert.NoError(t, err)

	chewable := &Chewable{}

	err = json.Unmarshal(dataRaw, chewable)
	assert.NoError(t, err)

	assert.EqualValues(t, &Chewable{
		Global: map[string]interface{}{
			"global_var":    float64(2),
			"overwrite_var": "global",
			"version":       float64(1),
		},
		Data: []ChewableData{
			{
				Templates: map[string]string{
					"t1":"t1.out",
					"t2":"t2.out",
				},
				Local: map[string]interface{}{
					"overwrite_var": "local",
					"local_var":     float64(3),
				},
			},
		},
	}, chewable)
}
