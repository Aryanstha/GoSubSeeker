package lib

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"os"
	"reflect"
	"strings"
)

type Options struct {
	Threads        int
	Domain         string
	Dict           string
	Depth          int
	Help           bool
	Log            string
	DNSServer      string
	WildcardDomain bool
	AXFC           bool
	ScanListFN     string
	ScanDomainList []string
}

func New() *Options {
	return &Options{}
}

func (opts *Options) existsDomain() bool {
	opts.ScanDomainList = []string{}
	for {
		if opts.ScanListFN != "" {
			f, err := os.Open(opts.ScanListFN)
			if err != nil {
				break
			}

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if s := strings.TrimSpace(scanner.Text()); s != "" {
					opts.ScanDomainList = append(opts.ScanDomainList, s)
				}
			}
			f.Close()
		}

		break
	}

	if len(opts.ScanDomainList) > 0 {
		return true
	}
	if opts.Domain != "" {
		opts.ScanDomainList = append(opts.ScanDomainList, opts.Domain)
		return true
	}

	return false
}

func (opts *Options) Validate() *multierror.Error {
	if opts.Help {
		flag.Usage()
		os.Exit(0)
	}

	var errorList *multierror.Error
	if !opts.existsDomain() {
		errorList = multierror.Append(errorList, fmt.Errorf("Domain (-d): Must be specified"))
	}

	if opts.Threads <= 0 {
		errorList = multierror.Append(errorList, fmt.Errorf("-t best > 0"))
	}
	if opts.Depth <= 0 {
		errorList = multierror.Append(errorList, fmt.Errorf("Depth scan (-depth): range [>=1]"))
	}

	_, err := os.Stat(opts.Dict)
	if err != nil {
		errorList = multierror.Append(errorList, fmt.Errorf("Dictionary file  (-f): Must be specified"))
	}
	if opts.Log == "" {
		logDir := "log"
		_, err := os.Stat(logDir)
		if err != nil {
			os.Mkdir(logDir, os.ModePerm)
		}
		opts.Log = fmt.Sprintf("%s/%s.txt", logDir, opts.Domain)
	}

	if opts.DNSServer == "" {
		opts.DNSServer = "8.8.8.8/8.8.4.4"
	}

	return errorList
}

func (opts *Options) PrintOptions() {
	value := reflect.ValueOf(*opts)
	types := reflect.TypeOf(*opts)

	fmt.Fprintln(os.Stderr, `=============================================
subdomain-scanner v0.4#dev
=============================================`)

	for i := 0; i < types.NumField(); i++ {
		if types.Field(i).Name[0] >= 65 && types.Field(i).Name[0] <= 90 {
			if value.Field(i).Interface() != "" {
				fmt.Fprintf(os.Stderr, "[+] %-15s: %v\n", types.Field(i).Name, value.Field(i).Interface())
			}
		}
	}
	fmt.Fprintln(os.Stderr, "=============================================")
}
