package log

type Line interface {
	LineKind() string
}

type Log struct {
	Kind          string
	Line          string
	Time          string
	Elb           string
	ClientIP      string
	ClientPort    string
	TargetIP      string
	TargetPort    string
	ReceivedBytes string
	SentBytes     string
	DomainName    string
}

type NLBLog struct {
	Log
	Type                     string
	Version                  string
	ListenerID               string
	TcpConnectionTimeMs      string
	TlsHandshakeTimeMs       string
	IncomingTlsAlert         string
	CertArn                  string
	CertificateSerial        string
	TlsCipherSuite           string
	TlsProtocolVersion       string
	TlsNamedGroup            string
	AlpnFeProtocol           string
	AlpnBeProtocol           string
	AlpnClientPreferenceList string
}

func (l NLBLog) LineKind() string {
	return "nlb"
}

type ALBLog struct {
	Log
	RequestProcessingTime  string
	TargetProcessingTime   string
	ResponseProcessingTime string
	ElbStatusCode          string
	TargetStatusCode       string
	RequestVerb            string
	RequestUrl             string
	RequestProto           string
	UserAgent              string
	SslCipher              string
	SslProtocol            string
	TargetGroupArn         string
	TraceID                string
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

func (l ALBLog) LineKind() string {
	return "alb"
}
