package parse

import (
	"regexp"

	"github.com/eminugurkenar/tailb/log"
)

type ALBParser struct {
	re *regexp.Regexp
}

func NewALBParser() *ALBParser {
	re := regexp.MustCompile(`(?P<type>[^ ]*) (?P<time>[^ ]*) (?P<elb>[^ ]*) (?P<client_ip>[^ ]*):(?P<client_port>[0-9]*) (?P<target_ip>[^ ]*)[:-](?P<target_port>[0-9]*) (?P<request_processing_time>[-.0-9]*) (?P<target_processing_time>[-.0-9]*) (?P<response_processing_time>[-.0-9]*) (?P<elb_status_code>|[-0-9]*) (?P<target_status_code>-|[-0-9]*) (?P<received_bytes>[-0-9]*) (?P<sent_bytes>[-0-9]*) \"(?P<request_verb>[^ ]*) (?P<request_url>.*) (?P<request_proto>- |[^ ]*)\" \"(?P<user_agent>[^\"]*)\" (?P<ssl_cipher>[A-Z0-9-_]+) (?P<ssl_protocol>[A-Za-z0-9.-]*) (?P<target_group_arn>[^ ]*) \"(?P<trace_id>[^\"]*)\" \"(?P<domain_name>[^\"]*)\" \"(?P<chosen_cert_arn>[^\"]*)\" (?P<matched_rule_priority>[-.0-9]*) (?P<request_creation_time>[^ ]*) \"(?P<actions_executed>[^\"]*)\" \"(?P<redirect_url>[^\"]*)\" \"(?P<lambda_error_reason>[^ ]*)\" \"(?P<target_port_list>[^\s]+?)\" \"(?P<target_status_code_list>[^\s]+)\" \"(?P<classification>[^ ]*)\" \"(?P<classification_reason>[^ ]*)\"`)
	return &ALBParser{
		re: re,
	}
}

func (p *ALBParser) Parse(line string) (log.Log, error) {
	l := log.Log{}
	match := p.re.FindAllStringSubmatch(line, -1)

	fields := p.re.SubexpNames()
	for i, m := range match[0] {
		switch fields[i] {
		case "type":
			l.Kind = m
		case "time":
			l.Time = m
		case "elb":
			l.Elb = m
		case "client_ip":
			l.ClientIP = m
		case "client_port":
			l.ClientPort = m
		case "target_ip":
			l.TargetIP = m
		case "target_port":
			l.TargetPort = m
		case "request_processing_time":
			l.RequestProcessingTime = m
		case "target_processing_time":
			l.TargetProcessingTime = m
		case "response_processing_time":
			l.ResponseProcessingTime = m
		case "elb_status_code":
			l.ElbStatusCode = m
		case "target_status_code":
			l.TargetStatusCode = m
		case "received_bytes":
			l.ReceivedBytes = m
		case "sent_bytes":
			l.SentBytes = m
		case "request_verb":
			l.RequestVerb = m
		case "request_url":
			l.RequestUrl = m
		case "request_proto":
			l.RequestProto = m
		case "user_agent":
			l.UserAgent = m
		case "ssl_cipher":
			l.SslCipher = m
		case "ssl_protocol":
			l.SslProtocol = m
		case "target_group_arn":
			l.TargetGroupArn = m
		case "trace_id":
			l.TraceID = m
		case "domain_name":
			l.DomainName = m
		case "chosen_cert_arn":
			l.ChosenCertArn = m
		case "matched_rule_priority":
			l.MatchedRulePriority = m
		case "request_creation_time":
			l.RequestCreationTime = m
		case "actions_executed":
			l.ActionsExecuted = m
		case "redirect_url":
			l.RedirectUrl = m
		case "lambda_error_reason":
			l.LambdaErrorReason = m
		case "target_port_list":
			l.TargetPortList = m
		case "target_status_code_list":
			l.TargetStatusCodeList = m
		case "classification":
			l.Classification = m
		case "classification_reason":
			l.ClassificationReason = m
		}

	}
	return l, nil
}
