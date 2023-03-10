package loadbalancer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type LoadBalancerFinder struct {
	client *elbv2.ELBV2
}

func NewLoadBalancerFinder(client *elbv2.ELBV2) *LoadBalancerFinder {
	return &LoadBalancerFinder{
		client: client,
	}

}

func (f *LoadBalancerFinder) ListLoadBalancers(ctx context.Context) ([]*LoadBalancer, error) {
	return f.listLoadBalancers(ctx, &elbv2.DescribeLoadBalancersInput{})
}

func (f *LoadBalancerFinder) GetLoadBalancer(ctx context.Context, name string) (*LoadBalancer, error) {
	lbs, err := f.listLoadBalancers(ctx, &elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String(name),
		},
	})
	if err != nil {
		return nil, err
	}
	return lbs[0], err
}

func (f *LoadBalancerFinder) listLoadBalancers(ctx context.Context, input *elbv2.DescribeLoadBalancersInput) ([]*LoadBalancer, error) {
	var lbs []*LoadBalancer

	var err error

	err = f.client.DescribeLoadBalancersPagesWithContext(ctx, input,
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			for _, lbo := range page.LoadBalancers {
				lb := LoadBalancer{kind: *lbo.Type, name: *lbo.LoadBalancerName, arn: *lbo.LoadBalancerArn}
				err = f.fillAccessLogInfo(ctx, &lb)
				if err != nil {
					return false
				}
				lbs = append(lbs, &lb)
			}
			return lastPage
		})

	if err != nil {
		return nil, err
	}

	return lbs, err
}

func (f *LoadBalancerFinder) fillAccessLogInfo(ctx context.Context, lb *LoadBalancer) error {
	input := &elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: aws.String(lb.arn),
	}

	result, err := f.client.DescribeLoadBalancerAttributesWithContext(ctx, input)
	if err != nil {
		return err
	}

	for _, attr := range result.Attributes {
		if *attr.Key == "access_logs.s3.bucket" {
			lb.accessLogBucket = *attr.Value
		}
		if *attr.Key == "access_logs.s3.prefix" {
			lb.accessLogBucketPrefix = *attr.Value
		}
	}

	sArn, err := arn.Parse(lb.arn)
	if err != nil {
		return err
	}

	lb.region = sArn.Region
	lb.accountId = sArn.AccountID
	lb.resource = sArn.Resource

	return nil
}
