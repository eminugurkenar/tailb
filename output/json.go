package output

import (
	"encoding/json"
	"fmt"

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

	var fieldValues map[string]interface{}

	switch kind {
	case "application":
		l := l.Line.(log.ALBLog)
		fieldValues = getFieldValues(l)
	case "network":
		l := l.Line.(log.NLBLog)
		fieldValues = getFieldValues(l)
	case "classic":
		l := l.Line.(log.ELBLog)
		fieldValues = getFieldValues(l)
	default:
		return nil, fmt.Errorf("Unknown log type %s", kind)
	}

	formatedLog := make(map[string]interface{})

	for k, v := range fieldValues {
		if !l.fieldFilter[k] {
			continue
		}

		formatedLog[k] = v
	}

	return json.Marshal(formatedLog)
}
