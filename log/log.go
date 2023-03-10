package log

type Log struct {
	Line                   string
	Kind                   string
	Time                   string
	Elb                    string
	ClientIP               string
	ClientPort             string
	TargetIP               string
	TargetPort             string
	RequestProcessingTime  string
	TargetProcessingTime   string
	ResponseProcessingTime string
	ElbStatusCode          string
	TargetStatusCode       string
	ReceivedBytes          string
	SentBytes              string
	RequestVerb            string
	RequestUrl             string
	RequestProto           string
	UserAgent              string
	SslCipher              string
	SslProtocol            string
	TargetGroupArn         string
	TraceID                string
	DomainName             string
	ChosenCertArn          string
	MatchedRulePriority    string
	RequestCreationTime    string
	ActionsExecuted        string
	RedirectUrl            string
	LambdaErrorReason      string
	TargetPortList         string
	TargetStatusCodeList   string
	Classification         string
	ClassificationReason   string
}
