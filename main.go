package main

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/eminugurkenar/tailb/between"
	"github.com/eminugurkenar/tailb/loadbalancer"
	"github.com/eminugurkenar/tailb/output"
	"github.com/eminugurkenar/tailb/parse"
	"github.com/eminugurkenar/tailb/read"
)

func printTable(lbs []*loadbalancer.LoadBalancer) {
	var maxNameL int
	var maxTypeL int
	var maxArnResourceL int

	for _, lb := range lbs {
		if len(lb.GetName()) > maxNameL {
			maxNameL = len(lb.GetName())
		}

		if len(lb.GetType()) > maxTypeL {
			maxTypeL = len(lb.GetName())
		}

		if len(lb.GetArnResource()) > maxArnResourceL {
			maxArnResourceL = len(lb.GetArnResource())
		}
	}

	maxNameL += 3
	maxTypeL += 3
	maxArnResourceL += 3

	for i, lb := range lbs {
		if i == 0 { // table header
			fmt.Printf("%-*s ", maxTypeL, "Type")
			fmt.Printf("%-*s ", maxNameL, "Name")
			fmt.Printf("%-*s ", maxArnResourceL, "Arn")
			fmt.Printf("%s \n", "Log Storage")
		}
		fmt.Printf("%-*s ", maxTypeL, lb.GetType())
		fmt.Printf("%-*s ", maxNameL, lb.GetName())
		fmt.Printf("%-*s ", maxArnResourceL, lb.GetArnResource())
		fmt.Printf("%s \n", lb.GetLogStorage())
	}
}

func listLoadBalancers(ctx context.Context, s *session.Session) error {
	elbc := elbv2.New(s)

	lbf := loadbalancer.NewLoadBalancerFinder(elbc)

	lbs, err := lbf.ListLoadBalancers(ctx)

	if err != nil {
		return err
	}

	printTable(lbs)

	return nil
}

func tail(ctx context.Context, s *session.Session, lbn string, start, end time.Time) error {
	elbc := elbv2.New(s)
	s3c := s3.New(s)

	lbf := loadbalancer.NewLoadBalancerFinder(elbc)
	lb, err := lbf.GetLoadBalancer(ctx, lbn)

	if err != nil {
		return err
	}

	b := between.NewBetween(start, end)
	days := b.ListDays()

	var prefixes []string

	for _, d := range days {
		prefixes = append(prefixes, fmt.Sprintf("%s/%s", lb.GetAccessLogPrefix(), d.Format("2006/01/02")))
	}

	p := parse.NewParser(lb.GetType())

	o := output.NewOutputFormater("json", output.DefaultOutputFormatOptions())

	lr := read.NewObjectStorageReader(s3c, lb.GetAccessLogBucket(), prefixes, ctx)

	scanner := bufio.NewScanner(lr)
	for scanner.Scan() {
		log, err := p.Parse(scanner.Text())
		if err != nil {
			return err
		}

		flog, err := o.Format(log)
		if err != nil {
			return err
		}
		fmt.Println(string(flog))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type AnyDate string

func (a *AnyDate) Validate() error {
	_, err := dateparse.ParseLocal(string(*a))
	return err
}

var CLI struct {
	Tail struct {
		Loadbalancer string  `arg:"" name:"loadbalancer" help:"Loadbalancer to tail." type:"string"`
		Start        AnyDate `help:"start date." short:"s" default:"${start}" required`
		End          AnyDate `help:"end date." short:"e" default:"${end}" required`
		Lines        int64   `help:"number of lines to print, (default -1, till end)" name:"nlines" short:"n" default:"-1" required`
	} `cmd:"" name:"tail" default:"withargs" help:"Tail logs of given loadbalancer. (default command)"`

	ListLoadBalancers struct {
	} `cmd:"" name:"ls" help:"List loadbalancers."`
}

func main() {
	s := session.New()
	ctx := context.Background()

	kctx := kong.Parse(&CLI,
		kong.Vars{"start": time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")},
		kong.Vars{"end": time.Now().Format("2006-01-02 15:04:05")},
	)

	switch kctx.Command() {
	case "tail <loadbalancer>":
		//already validated dates
		start, _ := dateparse.ParseLocal(string(CLI.Tail.Start))
		end, _ := dateparse.ParseLocal(string(CLI.Tail.End))

		if err := tail(ctx, s, CLI.Tail.Loadbalancer, start, end); err != nil {
			fmt.Println(err)
		}
	case "ls":
		if err := listLoadBalancers(ctx, s); err != nil {
			fmt.Println(err)
		}
	default:
		panic(kctx.Command())
	}
}
