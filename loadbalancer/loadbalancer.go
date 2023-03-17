package loadbalancer

import "fmt"

type LoadBalancer struct {
	kind                  string
	arn                   string
	name                  string
	accessLogBucket       string
	accessLogBucketPrefix string
	accountId             string
	region                string
	resource              string
}

func (l *LoadBalancer) GetAccessLogPrefix() string {
	if l.accessLogBucket == "" {
		return ""

	}

	return fmt.Sprintf("%s/AWSLogs/%s/elasticloadbalancing/%s", l.accessLogBucketPrefix, l.accountId, l.region)
}

func (l *LoadBalancer) GetName() string {
	return l.name
}

func (l *LoadBalancer) GetAccessLogBucket() string {
	return l.accessLogBucket
}

func (l *LoadBalancer) GetArn() string {
	return l.arn
}

func (l *LoadBalancer) GetKind() string {
	return l.kind
}

func (l *LoadBalancer) GetAccessLogBucketPrefix() string {
	return l.accessLogBucketPrefix
}

func (l *LoadBalancer) GetArnResource() string {
	return l.resource
}

func (l *LoadBalancer) String() string {
	return fmt.Sprintf("%s %s", l.name, l.arn)
}

func (l *LoadBalancer) GetLogStorage() string {
	if l.accessLogBucket == "" {
		return ""

	}

	return fmt.Sprintf("%s/%s", l.accessLogBucket, l.accessLogBucketPrefix)
}
