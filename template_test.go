package chew

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmth(t *testing.T) {
	dataRaw, err := ioutil.ReadFile("plsql/data/data.json")
	assert.NoError(t, err)

	chewable := &Chewable{}

	err = json.Unmarshal(dataRaw, chewable)
	assert.NoError(t, err)

	template := New("main")
	_, err = template.ParseFolder("plsql/templates")
	assert.NoError(t, err)

	template.ExecuteChewable(WriterWrapper{os.Stdout}, *chewable)
}
