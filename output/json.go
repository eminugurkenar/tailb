package output

import (
	"encoding/json"
	"reflect"

	"github.com/eminugurkenar/tailb/log"
)

type JsonLog struct {
	log.Log
	fieldFilter map[string]bool
}

type JSONOutputFormater struct {
	fieldFilter map[string]bool
}

func NewJSONOutputFormater(o OutputFormatOptions) OutputFormater {
	return &JSONOutputFormater{
		fieldFilter: createFieldFilter(o.Fields),
	}
}

func (f *JSONOutputFormater) Format(l log.Log) ([]byte, error) {
	return json.Marshal(JsonLog{l, f.fieldFilter})
}

func (l JsonLog) MarshalJSON() ([]byte, error) {
	rt, rv := reflect.TypeOf(l.Log), reflect.ValueOf(l.Log)

	formatedLog := make(map[string]interface{})

	for i := 0; i < rt.NumField(); i++ {
		if !l.fieldFilter[rt.Field(i).Name] {
			continue
		}

		formatedLog[rt.Field(i).Name] = rv.Field(i).Interface()
	}

	return json.Marshal(formatedLog)
}

func createFieldFilter(fields []string) map[string]bool {
	filter := make(map[string]bool)

	for _, field := range fields {
		filter[field] = true
	}

	return filter
}
