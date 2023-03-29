package parse

import (
	"regexp"

	"github.com/eminugurkenar/tailb/log"
)

type ELBParser struct {
	re *regexp.Regexp
}

func NewELBParser() *ELBParser {
	re := regexp.MustCompile(`(?P<timestamp>[^ ]*) (?P<elb_name>[^ ]*) (?P<request_ip>[^ ]*):(?P<request_port>[0-9]*) (?P<backend_ip>[^ ]*)[:-](?P<backend_port>[0-9]*) (?P<request_processing_time>[-.0-9]*) (?P<backend_processing_time>[-.0-9]*) (?P<client_response_time>[-.0-9]*) (?P<elb_response_code>[-0-9]*) (?P<backend_response_code>[-0-9]*) (?P<received_bytes>[|[-0-9]*) (?P<sent_bytes>[-|[-0-9]*) (?P<request_verb>[-0-9]*) (?P<url>[^ ]*) (?P<protocol>(- |[^ ]*)) (?P<user_agent>\"[^\"]*\") (?P<ssl_cipher>[A-Z0-9-]+) (?P<ssl_protocol>[A-Za-z0-9.-]*)$`)
	return &ELBParser{
		re: re,
	}
}

func (p *ELBParser) Parse(line string) (log.Line, error) {
	l := log.ELBLog{}

	match := p.re.FindStringSubmatch(line)
	fields := p.re.SubexpNames()

	for i, m := range match {
		switch fields[i] {
		case "timestamp":
			l.Time = m
		case "elb_name":
			l.Elb = m
		case "request_ip":
			l.ClientIP = m
		case "request_port":
			l.ClientPort = m
		case "backend_ip":
			l.TargetIP = m
		case "backend_port":
			l.TargetPort = m
		case "request_processing_time":
			l.RequestProcessingTime = m
		case "backend_processing_time":
			l.BackendProcessingTime = m
		case "client_response_time":
			l.ClientResponseTime = m
		case "elb_response_code":
			l.ElbResponseCode = m
		case "backend_response_code":
			l.BackendResponseCode = m
		case "received_bytes":
			l.ReceivedBytes = m
		case "sent_bytes":
			l.SentBytes = m
		case "request_verb":
			l.RequestVerb = m
		case "url":
			l.URL = m
		case "protocol":
			l.Protocol = m
		case "user_agent":
			l.UserAgent = m
		case "ssl_cipher":
			l.SSLCipher = m
		case "ssl_protocol":
			l.SSLProtocol = m
		}
	}

	return l, nil

}
