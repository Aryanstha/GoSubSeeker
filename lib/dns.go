package lib

import (
	"github.com/aryanstha/GoSubSeeker/lib/dns"
	"net"
	"strings"
)

func (this *Scanner) LookupHost(host string) (addrs []net.IP, err error) {
	DnsResolver := dns.New(strings.Split(this.opts.DNSServer, "/"))

	ipaddr, err := DnsResolver.LookupHost(host)
	if err != nil {
		//log.Println(err)
		return
	}

	return ipaddr, nil
}

func (this *Scanner) LookupNS(host string) ([]string, error) {
	DnsResolver := dns.New(strings.Split(this.opts.DNSServer, "/"))

	return DnsResolver.LookupNS(host)
}
