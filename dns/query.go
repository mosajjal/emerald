package dns

import (
	"context"
	"net"
	"time"
)

// QueryTXT function takes a FQDN and a server as input argument and
// returns the TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
// QueryTXT does not perform any checks on the FQDN's validity. the caller
// should take care of that
func QueryTXT(ctx context.Context, fqdn string, server net.IP) (output string, err error) {
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

	if r, e := resolver.LookupTXT(ctx, fqdn); e == nil {
		output = r[0]
	} else {
		err = e
	}
	return
}
