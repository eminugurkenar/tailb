package parse

import "github.com/eminugurkenar/tailb/log"

type Parser interface {
	Parse(line string) (log.Log, error)
}

func NewParser(kind string) Parser {
	return NewALBParser()
}
