package output

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/eminugurkenar/tailb/log"
)

type JsonLog struct {
	log.Line
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

func (f *JSONOutputFormater) Format(l log.Line) ([]byte, error) {
	return json.Marshal(JsonLog{l, f.fieldFilter})
}

func (l JsonLog) MarshalJSON() ([]byte, error) {
	kind := l.LineKind()

	var rtb reflect.Type
	var rvb reflect.Value

	var rtl reflect.Type
	var rvl reflect.Value

	switch kind {
	case "alb":
		l := l.Line.(log.ALBLog)
		rtb, rvb = reflect.TypeOf(l.Log), reflect.ValueOf(l.Log)
		rtl, rvl = reflect.TypeOf(l), reflect.ValueOf(l)
	case "nlb":
		l := l.Line.(log.NLBLog)
		rtb, rvb = reflect.TypeOf(l.Log), reflect.ValueOf(l.Log)
		rtl, rvl = reflect.TypeOf(l), reflect.ValueOf(l)
	default:
		return nil, fmt.Errorf("Unknown log type %s", kind)
	}

	formatedLog := make(map[string]interface{})

	for i := 0; i < rtb.NumField(); i++ {
		if !l.fieldFilter[rtb.Field(i).Name] {
			continue
		}

		formatedLog[rtb.Field(i).Name] = rvb.Field(i).Interface()
	}

	for i := 0; i < rtl.NumField(); i++ {
		if !l.fieldFilter[rtl.Field(i).Name] {
			continue
		}

		formatedLog[rtl.Field(i).Name] = rvl.Field(i).Interface()
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
