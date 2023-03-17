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
	var maxKindL int
	var maxArnResourceL int

	for _, lb := range lbs {
		if len(lb.GetName()) > maxNameL {
			maxNameL = len(lb.GetName())
		}

		if len(lb.GetKind()) > maxKindL {
			maxKindL = len(lb.GetName())
		}

		if len(lb.GetArnResource()) > maxArnResourceL {
			maxArnResourceL = len(lb.GetArnResource())
		}
	}

	maxNameL += 3
	maxKindL += 3
	maxArnResourceL += 3

	for i, lb := range lbs {
		if i == 0 { // table header
			fmt.Printf("%-*s ", maxKindL, "Type")
			fmt.Printf("%-*s ", maxNameL, "Name")
			fmt.Printf("%-*s ", maxArnResourceL, "Arn")
			fmt.Printf("%s \n", "Log Storage")
		}
		fmt.Printf("%-*s ", maxKindL, lb.GetKind())
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

func tail(ctx context.Context, s *session.Session, lbn string, start, end time.Time, fields []string) error {
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

	p, err := parse.NewParser(lb.GetKind())
	if err != nil {
		return err
	}

	of, err := output.DefaultOutputFormatOptions(lb.GetKind())
	if err != nil {
		return err
	}

	if len(fields) > 0 {
		of, err = output.NewOutputFormatOptions(lb.GetKind(), fields)
	}

	o := output.NewOutputFormater("json", of)

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
		Loadbalancer string   `arg:"" name:"loadbalancer" help:"Loadbalancer to tail."`
		Start        AnyDate  `short:"s" default:"${start}" required help:"start date."`
		End          AnyDate  `short:"e" help:"end date." default:"${end}." required`
		Lines        int64    `short:"n" default:"-1" required help:"number of lines to print, (default -1, till end)." name:"nlines"`
		Fields       []string `short:"f" help:"fields to be shown in output."`
	} `cmd:"" name:"tail" default:"withargs" help:"tail logs of given loadbalancer (default command)."`

	ListLoadBalancers struct {
	} `cmd:"" name:"ls" help:"list loadbalancers."`
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

		if err := tail(ctx, s, CLI.Tail.Loadbalancer, start, end, CLI.Tail.Fields); err != nil {
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
