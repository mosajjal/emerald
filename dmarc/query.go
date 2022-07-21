package dmarc

import (
	"context"
	"net"
	"strings"
	"time"
)

// Query function takes a top level domain name (google.com) and
// returns the DMARC TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
func Query(ctx context.Context, domain string, server net.IP) (output string, err error) {
	if !strings.HasPrefix(domain, "_dmarc.") {
		domain = "_dmarc." + domain
	}
	domain = strings.TrimSuffix(domain, ".")

	resolver := net.DefaultResolver

	if !server.Equal(net.IPv4zero) {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(10000),
				}
				return d.DialContext(ctx, network, address+":53")
			},
		}
	}

	if r, e := resolver.LookupTXT(ctx, domain); e == nil {
		output = r[0]
	} else {
		err = e
	}
	return
}
