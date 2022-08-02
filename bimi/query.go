package bimi

import (
	"context"
	"net"
	"strings"

	"github.com/mosajjal/emerald/dns"
)

// Query function takes a top level domain name (google.com) and
// returns the bimi TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
func Query(ctx context.Context, domain string, server net.IP) (output string, err error) {
	if !strings.HasPrefix(domain, "_bimi.") {
		domain = "default._bimi." + domain
	}
	domain = strings.TrimSuffix(domain, ".")

	return dns.QueryTXT(ctx, domain, server)

}
