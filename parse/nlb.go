package parse

import (
	"fmt"
	"regexp"

	"github.com/eminugurkenar/tailb/log"
)

type NLBParser struct {
	re *regexp.Regexp
}

func NewNLBParser() *NLBParser {
	return &NLBParser{
		re: nil,
	}
}

func (p *NLBParser) Parse(line string) (log.Log, error) {
	return log.Log{}, fmt.Errorf("NLB log parsing not implemented")
}
