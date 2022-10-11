package dns

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"

	mkdns "github.com/miekg/dns"
)

type DnsClient struct {
	C   net.Conn
	ctx context.Context
}

func (d DnsClient) QueryMX(ctx context.Context, fqdn string) ([]mkdns.RR, error) {
	m := new(mkdns.Msg)
	if !strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn + "."
	}
	m.SetQuestion(fqdn, mkdns.TypeMX)
	return d.performQuery(m)
}

func (d DnsClient) QueryTXT(ctx context.Context, fqdn string) ([]mkdns.RR, error) {
	m := new(mkdns.Msg)
	if !strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn + "."
	}
	m.SetQuestion(fqdn, mkdns.TypeTXT)
	return d.performQuery(m)
}

func (d DnsClient) performQuery(q *mkdns.Msg) (responses []mkdns.RR, err error) {

	fnDone := make(chan bool)
	go func() {
		co := &mkdns.Conn{Conn: d.C}
		if err = co.WriteMsg(q); err != nil {
			fnDone <- true
		}
		var r *mkdns.Msg
		r, err = co.ReadMsg()
		co.Close()
		if err == nil {
			if r.Truncated {
				err = fmt.Errorf("response was truncated. consider using a different protocol (TCP) for large queries")
			} else if r.Id == q.Id {
				responses = r.Answer
			} else {
				err = fmt.Errorf("%d", r.Id)
			}
		}
		fnDone <- true
	}()
	for {
		select {
		case <-fnDone:
			return
		case <-d.ctx.Done():
			return
		}
	}
}

func NewDnsClient(ctx context.Context, server string) (DnsClient, error) {
	dnsUrl, err := url.Parse(server)
	c := DnsClient{ctx: ctx}
	//TODO: none of the below connections have any context so timeout doesn't apply to them

	if err == nil {
		switch dnsUrl.Scheme {
		case "udp", "udp4":
			s, err := net.ResolveUDPAddr("udp4", dnsUrl.Host)
			if err == nil {
				c.C, err = net.DialUDP("udp4", nil, s)
				return c, err
			}
			return c, err
		case "udp6":
			s, err := net.ResolveUDPAddr("udp6", dnsUrl.Host)
			if err == nil {
				c.C, err = net.DialUDP("udp6", nil, s)
				return c, err
			}
			return c, err
		case "tcp", "tcp4":
			s, err := net.ResolveTCPAddr("tcp4", dnsUrl.Host)
			if err == nil {
				c.C, err = net.DialTCP("tcp4", nil, s)
				return c, err
			}
			return c, err
		case "tcp6":
			s, err := net.ResolveTCPAddr("tcp6", dnsUrl.Host)
			if err == nil {
				c.C, err = net.DialTCP("tcp6", nil, s)
				return c, err
			}
			return c, err
		case "tls", "tls4":
			// s, err := net.ResolveTCPAddr("tcp4", dnsUrl.Host)
			if err == nil {
				c.C, err = tls.Dial("tcp4", dnsUrl.Host, &tls.Config{})
				return c, err
			}
		case "tls6":
			// s, err := net.ResolveTCPAddr("tcp4", dnsUrl.Host)
			if err == nil {
				c.C, err = tls.Dial("tcp6", dnsUrl.Host, &tls.Config{})
				return c, err
			}
		}
		//TODO: build DoH, DoQ, here

	} else {
		return c, err
	}
	return c, nil
}
