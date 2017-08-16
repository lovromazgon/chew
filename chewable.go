package chew

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Chewable struct {
	Global map[string]interface{}
	Data   []ChewableData
}

type ChewableData struct {
	Templates map[string]string
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
	for i, d := range dataSlice {
		cd, err := extractChewableData(d)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not extract object %d in field 'data': %v", i, err))
		}
		c.Data[i] = cd
	}

	return nil
}

func extractChewableData(data interface{}) (cd ChewableData, err error) {
	local, err := ToMap(data)
	if err != nil {
		return cd, err
	}

	templatesRaw, ok := local["templates"]
	if !ok {
		return cd, errors.New("Could not find field 'templates'")
	}
	templatesMap, err := ToMap(templatesRaw)
	if err != nil {
		return cd, err
	}

	templates := make(map[string]string)
	for tmpl, outRaw := range templatesMap {
		outString, ok := outRaw.(string)
		if !ok {
			return cd, errors.New(fmt.Sprintf("Value of %s in 'templates' is not a string", tmpl))
		}
		templates[tmpl] = outString
	}

	delete(local, "templates")
	cd.Local = local
	cd.Templates = templates

	return
}

func ToMap(data interface{}) (map[string]interface{}, error) {
	var dataMap map[string]interface{}
	var ok bool

	if dataMap, ok = data.(map[string]interface{}); ok {
		// nothing to do, this is already a map
	} else {
		dataVal := reflect.Indirect(reflect.ValueOf(data))
		if dataVal.Kind() == reflect.Struct {
			dataMap = make(map[string]interface{})

			for i := 0; i < dataVal.NumField(); i++ {
				fieldName := dataVal.Type().Field(i).Name
				fieldValue := dataVal.Field(i).Interface()
				dataMap[fieldName] = fieldValue
			}
		}
	}

	if dataMap == nil {
		return nil, errors.New("Could not extract map!")
	}

	return dataMap, nil
}