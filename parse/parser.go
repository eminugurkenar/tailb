package parse

import (
	"fmt"

	"github.com/eminugurkenar/tailb/log"
)

type Parser interface {
	Parse(line string) (log.Log, error)
}

func NewParser(kind string) (Parser, error) {
	switch kind {
	case "application":
		return NewALBParser(), nil
	case "network":
		return NewNLBParser(), nil
	}
	return nil, fmt.Errorf("parser not available for %s", kind)
}
