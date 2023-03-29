package output

import (
	"fmt"
	"reflect"

	"github.com/eminugurkenar/tailb/log"
)

type OutputFormater interface {
	Format(l log.Line) ([]byte, error)
}

type OutputFormatOptions struct {
	Fields []string
}

func DefaultOutputFormatOptions(kind string) (OutputFormatOptions, error) {
	var fieldsValues map[string]interface{}
	switch kind {
	case "application":
		fieldsValues = getFieldValues(log.ALBLog{})
	case "network":
		fieldsValues = getFieldValues(log.NLBLog{})
	case "classic":
		fieldsValues = getFieldValues(log.ELBLog{})
	default:
		return OutputFormatOptions{}, fmt.Errorf("outputoptions not available for %s", kind)
	}

	fields := make([]string, len(fieldsValues))

	i := 0
	for k := range fieldsValues {
		fields[i] = k
		i++
	}

	return OutputFormatOptions{Fields: fields}, nil

}

func NewOutputFormatOptions(kind string, fields []string) (OutputFormatOptions, error) {
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	// Filter the default fields to only include those in the map
	filteredFields := []string{}
	options, err := DefaultOutputFormatOptions(kind)
	if err != nil {
		return options, err
	}

	for _, field := range options.Fields {
		if fieldMap[field] {
			filteredFields = append(filteredFields, field)
		}
	}

	return OutputFormatOptions{
		Fields: filteredFields,
	}, nil
}

func NewOutputFormater(kind string, o OutputFormatOptions) OutputFormater {
	return NewJSONOutputFormater(o)
}

func createFieldFilter(fields []string) map[string]bool {
	filter := make(map[string]bool)

	for _, field := range fields {
		filter[field] = true
	}

	return filter
}

func getFieldValues(v interface{}) map[string]interface{} {
	values := make(map[string]interface{})

	t := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := vv.Field(i)

		if field.Anonymous {
			embeddedValues := getFieldValues(value.Interface())
			for k, v := range embeddedValues {
				values[k] = v
			}
		} else {
			values[field.Name] = value.Interface()
		}
	}

	return values
}
