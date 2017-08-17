package chew

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Chewable is the main object which carries the data in chew. It stores everything that is used
// when generating the data. The slice of ChewableData is looped through and used in executing
// the specified templates. Chewable also carries a map of objects called Global, which is accessible
// in every template.
type Chewable struct {
	Global map[string]interface{}
	Data   []ChewableData
}

// ChewableData is a collection of data which can be used in executing one or more templates.
// The field Templates stores a map in which the key denotes the name of the template which will be generated
// and the value denotes the output filename. Everything in field Local will be accessible in the templates.
// If a key in the map Local also exists in the map Global in Chewable, the Local value will be used.
type ChewableData struct {
	Templates map[string]string
	Local     map[string]interface{}
}

// UnmarshalJSON parses data from JSON into Chewable. Returns an error if parsing was unsuccessful, else nil.
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
			return fmt.Errorf("Could not extract object %d in field 'data': %v", i, err)
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
			return cd, fmt.Errorf("Value of %s in 'templates' is not a string", tmpl)
		}
		templates[tmpl] = outString
	}

	delete(local, "templates")
	cd.Local = local
	cd.Templates = templates

	return
}

// ToMap takes an object and extracts a map[string]interface{}. If the object is already a map[string]interface{}
// it is returned as it is. If the object is a struct, then the fields of the struct are mapped to a map
// where the keys are the names of the fields, while the values are the actual values. If anything else is
// passed to the function an error is thrown.
func ToMap(data interface{}) (map[string]interface{}, error) {
	var dataMap map[string]interface{}
	var ok bool

	if dataMap, ok = data.(map[string]interface{}); ok {
		// nothing to do, this is already a map
	} else {
		dataVal := reflect.Indirect(reflect.ValueOf(data))
		switch dataVal.Kind() {
		case reflect.Struct:
			dataMap = make(map[string]interface{})

			for i := 0; i < dataVal.NumField(); i++ {
				fieldName := dataVal.Type().Field(i).Name
				fieldValue := dataVal.Field(i).Interface()
				dataMap[fieldName] = fieldValue
			}

		default:
			return nil, fmt.Errorf("Could not extract map from type %s", dataVal.Type())

		}
	}

	return dataMap, nil
}
