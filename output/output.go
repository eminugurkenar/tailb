package output

import "github.com/eminugurkenar/tailb/log"

type OutputFormater interface {
	Format(l log.Log) ([]byte, error)
}

type OutputFormatOptions struct {
	Fields []string
}

func DefaultOutputFormatOptions() OutputFormatOptions {
	return OutputFormatOptions{
		Fields: []string{
			"Time",
			"Elb",
			"ClientIP",
			"ClientPort",
			"TargetIP",
			"TargetPort",
			"RequestProcessingTime",
			"TargetProcessingTime",
			"ResponseProcessingTime",
			"ElbStatusCode",
			"TargetStatusCode",
			"ReceivedBytes",
			"SentBytes",
			"RequestVerb",
			"RequestUrl",
			"RequestProto",
			"UserAgent",
			"SslCipher",
			"SslProtocol",
			"TargetGroupArn",
			"TraceID",
			"DomainName",
			"ChosenCertArn",
			"MatchedRulePriority",
			"RequestCreationTime",
			"ActionsExecuted",
			"RedirectUrl",
			"LambdaErrorReason",
			"TargetPortList",
			"TargetStatusCodeList",
			"Classification",
			"ClassificationReason",
		},
	}
}

func NewOutputFormater(kind string, o OutputFormatOptions) OutputFormater {
	return NewJSONOutputFormater(o)
}
