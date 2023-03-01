package dns

import (
	"net"
	"strings"
	"errors"
	"github.com/miekg/dns"
)

func Axrf(hostname string, servers []string) (results []string, err error) {
	results = []string{}
	domain := strings.ToLower(hostname)

	fqdn := dns.Fqdn(domain)
	for _, server := range servers {
		results = []string{}

		msg := new(dns.Msg)
		msg.SetAxfr(fqdn)

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server, "53"))
		if err != nil {
			continue
		}

		for a := range answerChan {

			if a.Error != nil {
				continue
			}

			for _, rr := range a.RR {
				results = append(results, rr.String())
			}
		}

		if (len(results) == 0) {
			continue
		}
		return results, nil
	}
	return results, errors.New("Transfer failed")
}
