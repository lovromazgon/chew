package chew

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Chewable struct {
	Global map[string]interface{}
	Data   []ChewableData
}

type ChewableData struct {
	Templates []string
	Local     map[string]interface{}
}

func (c *Chewable) UnmarshalJSON(data []byte) error {
	global := make(map[string]interface{})

	if err := json.Unmarshal(data, &global); err != nil {
		return err
	}

	dataObj := global["data"]
	if dataObj == nil {
		return errors.New("Could not find field 'data'")
	}

	delete(global, "data")
	c.Global = global

	dataSlice, ok := dataObj.([]interface{})
	if !ok {
		return errors.New("Field 'data' is not a slice")
	}

	c.Data = make([]ChewableData, len(dataSlice))
	for i,d := range dataSlice {
		cd, err := extractChewableData(d)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not extract object %d in field 'data': %v", i, err))
		}
		c.Data[i] = cd
	}

	return nil
}

func extractChewableData(data interface{}) (cd ChewableData, err error) {
	local, ok := data.(map[string]interface{})
	if !ok {
		return cd, errors.New("Object is not of type map[string]interface{}")
	}

	templatesRaw, ok := local["templates"]
	if !ok {
		return cd, errors.New("Could not find field 'templates'")
	}
	templatesSlice, ok := templatesRaw.([]interface{})
	if !ok {
		return cd, errors.New("Field 'templates' is not a slice")
	}

	templates := make([]string, len(templatesSlice))
	for i,tmplRaw := range templatesSlice {
		tmplString, ok := tmplRaw.(string)
		if !ok {
			return cd, errors.New(fmt.Sprintf("Object %d in 'templates' is not a string", i))
		}
		templates[i] = tmplString
	}

	delete(local, "templates")
	cd.Local = local
	cd.Templates = templates

	return
}

/*
func (u *Chewable) MarshalJSON() ([]byte, error) {

}
*/
