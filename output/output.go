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

func NewOutputFormatOptions(fields []string) OutputFormatOptions {
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	// Filter the default fields to only include those in the map
	filteredFields := []string{}
	for _, field := range DefaultOutputFormatOptions().Fields {
		if fieldMap[field] {
			filteredFields = append(filteredFields, field)
		}
	}

	return OutputFormatOptions{
		Fields: filteredFields,
	}
}

func NewOutputFormater(kind string, o OutputFormatOptions) OutputFormater {
	return NewJSONOutputFormater(o)
}
