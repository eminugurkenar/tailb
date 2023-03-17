package parse

import (
	"regexp"

	"github.com/eminugurkenar/tailb/log"
)

type NLBParser struct {
	re *regexp.Regexp
}

func NewNLBParser() *NLBParser {
	re := regexp.MustCompile(`(?P<type>[^ ]*) (?P<version>[^ ]*) (?P<time>[^ ]*) (?P<elb>[^ ]*) (?P<listener_id>[^ ]*) (?P<client_ip>[^ ]*):(?P<client_port>[0-9]*) (?P<target_ip>[^ ]*):(?P<target_port>[0-9]*) (?P<tcp_connection_time_ms>[-.0-9]*) (?P<tls_handshake_time_ms>[-.0-9]*) (?P<received_bytes>[-0-9]*) (?P<sent_bytes>[-0-9]*) (?P<incoming_tls_alert>[-0-9]*) (?P<cert_arn>[^ ]*) (?P<certificate_serial>[^ ]*) (?P<tls_cipher_suite>[^ ]*) (?P<tls_protocol_version>[^ ]*) (?P<tls_named_group>[^ ]*) (?P<domain_name>[^ ]*) (?P<alpn_fe_protocol>[^ ]*) (?P<alpn_be_protocol>[^ ]*) (?P<alpn_client_preference_list>[^ ]*)$`)
	return &NLBParser{
		re: re,
	}
}

func (p *NLBParser) Parse(line string) (log.Line, error) {
	l := log.NLBLog{}

	match := p.re.FindStringSubmatch(line)
	fields := p.re.SubexpNames()

	for i, m := range match {
		switch fields[i] {
		case "type":
			l.Type = m
		case "version":
			l.Version = m
		case "time":
			l.Time = m
		case "elb":
			l.Elb = m
		case "listener_id":
			l.ListenerID = m
		case "client_ip":
			l.ClientIP = m
		case "client_port":
			l.ClientPort = m
		case "target_ip":
			l.TargetIP = m
		case "target_port":
			l.TargetPort = m
		case "tcp_connection_time_ms":
			l.TcpConnectionTimeMs = m
		case "tls_handshake_time_ms":
			l.TlsHandshakeTimeMs = m
		case "received_bytes":
			l.ReceivedBytes = m
		case "sent_bytes":
			l.SentBytes = m
		case "incoming_tls_alert":
			l.IncomingTlsAlert = m
		case "cert_arn":
			l.CertArn = m
		case "certificate_serial":
			l.CertificateSerial = m
		case "tls_cipher_suite":
			l.TlsCipherSuite = m
		case "tls_protocol_version":
			l.TlsProtocolVersion = m
		case "tls_named_group":
			l.TlsNamedGroup = m
		case "domain_name":
			l.DomainName = m
		case "alpn_fe_protocol":
			l.AlpnFeProtocol = m
		case "alpn_be_protocol":
			l.AlpnBeProtocol = m
		case "alpn_client_preference_list":
			l.AlpnClientPreferenceList = m
		}
	}

	return l, nil
}
