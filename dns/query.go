package dns

import (
	"context"
	"strings"
	"time"

	"github.com/mosajjal/dnsclient"

	mkdns "github.com/miekg/dns"
)

func QueryMX(ctx context.Context, client dnsclient.Client, fqdn string) ([]mkdns.RR, time.Duration, error) {
	m := new(mkdns.Msg)
	if !strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn + "."
	}
	m.SetQuestion(fqdn, mkdns.TypeMX)
	return client.Query(ctx, m)
}

func QueryTXT(ctx context.Context, client dnsclient.Client, fqdn string) ([]mkdns.RR, time.Duration, error) {
	m := new(mkdns.Msg)
	if !strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn + "."
	}
	m.SetQuestion(fqdn, mkdns.TypeTXT)
	return client.Query(ctx, m)
}
