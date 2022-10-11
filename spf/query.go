package spf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"

	mkdns "github.com/miekg/dns"
	"github.com/mosajjal/emerald/dns"
	"gopkg.in/yaml.v3"
)

type SpfRecord struct {
	QueryDomain string   `desc:"Domain"`
	Version     string   `desc:"SPF record version"`
	IPs         []net.IP `desc:"List of allowed IP addresses authorised to send email on your behalf"`
	Includes    []string `desc:"List of allowed third party domains that can send email on you behalf"`
	Tag         string   `desc:"Tag (policy) applied to SPF. -all: fail, ~all: softfail, and +all: allow"`
}

// Query function takes a top level domain name (google.com) and
// returns the MX TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
func (s *SpfRecord) Query(ctx context.Context, server string) (err error) {
	c, err := dns.NewDnsClient(ctx, server)
	if err != nil {
		return err
	}
	// since TXT response will be multiple lines, the one with v=spf* will be the SPF record.
	responses, err := c.QueryTXT(ctx, s.QueryDomain)
	if err != nil {
		return err
	}
	for _, r := range responses {
		if strings.Contains(r.String(), "v=spf") {
			if t, ok := r.(*mkdns.TXT); ok {
				spfParts := strings.Split(t.Txt[0], " ")
				for _, kv := range spfParts {

					// ?all and v=spf? do not follow the same rules of key:value in SPF. need to be treated differently.
					if strings.Contains(kv, "v=spf") {
						s.Version = strings.Split(kv, "=")[1]
					}
					if strings.Contains(kv, "all") {
						s.Tag = kv
					}

					// trim the whitespace
					kv = strings.TrimSpace(kv)
					tmp := strings.SplitN(kv, ":", 2)
					if len(tmp) != 2 {
						continue
					}
					key, value := tmp[0], tmp[1]
					switch key {
					case "ip4", "ip6":
						s.IPs = append(s.IPs, net.ParseIP(value))
					case "include":
						s.Includes = append(s.Includes, value)
					}
				}
			} else {
				return fmt.Errorf("Can't parse MX response")
			}
		}
	}
	return
}

// Marshal provides a way to show a DMARC report in different formats
func (s SpfRecord) Marshal(kind string) ([]byte, error) {
	switch kind {
	case "json":
		return json.Marshal(s)
	case "yaml":
		return yaml.Marshal(s)
	case "pretty":
		var b bytes.Buffer
		_ = io.Writer(&b)
		dns.PrettyPrint(s, &b, "desc")
		return ioutil.ReadAll(&b)
	case "STIX":
		return nil, fmt.Errorf("STIX has not been implemented yet")
	}
	return nil, fmt.Errorf("Unknown kind: %s", kind)
}
