package chew

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmth(t *testing.T) {
	dataRaw, err := ioutil.ReadFile("data/data.json")
	assert.NoError(t, err)

	chewable := &Chewable{}

	err = json.Unmarshal(dataRaw, chewable)
	assert.NoError(t, err)

	myChew := New("main")
	_, err = myChew.ParseFolder("templates/plsql")
	assert.NoError(t, err)

	myChew.ExecuteChewable(os.Stdout, *chewable)
}
