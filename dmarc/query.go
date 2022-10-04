package dmarc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	mkdns "github.com/miekg/dns"
	"github.com/mosajjal/emerald/dns"
	yaml "gopkg.in/yaml.v3"
)

type DmarcDns struct {
	V     string `desc:"Protocol version"`
	Pct   uint8  `desc:"Percentage of messages subjected to filtering"`
	Ruf   string `desc:"Reporting URI for forensic reports"`
	Rua   string `desc:"Reporting URI of aggregate reports"`
	P     string `desc:"Policy for organizational domain"`
	Sp    string `desc:"Policy for subdomains of the OD"`
	Adkim string `desc:"Alignment mode for DKIM"`
	Aspf  string `desc:"Alignment mode for SPF"`
}

// Query function takes a top level domain name (google.com) and
// returns the DMARC TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
func Query(ctx context.Context, domain string, server string) (d DmarcDns, err error) {
	if !strings.HasPrefix(domain, "_dmarc.") {
		domain = "_dmarc." + domain
	}
	domain = strings.TrimSuffix(domain, ".")
	c, err := dns.NewDnsClient(ctx, server)
	if err != nil {
		return
	}
	responses, err := c.QueryTXT(ctx, domain)
	if len(responses) == 0 {
		return d, fmt.Errorf("no DMARC response")
	}
	// [_dmarc.n0p.me.	300	IN	TXT	"v=DMARC1; p=quarantine; pct=100; rua=mailto:hi@n0p.me;"]

	if t, ok := responses[0].(*mkdns.TXT); ok {
		dmarcParts := strings.Split(t.Txt[0], ";")
		for _, kv := range dmarcParts {
			// trim the whitespace
			kv = strings.TrimSpace(kv)
			if len(kv) == 0 {
				continue
			}
			// we need to split kv by = to get k and v individually
			tmp := strings.SplitN(kv, "=", 2)
			if len(tmp) != 2 {
				return d, fmt.Errorf("Wrong TXT response: %s", kv)
			}
			key, value := tmp[0], tmp[1]

			switch key {
			case "v":
				d.V = value
			case "pct":
				pct, _ := strconv.Atoi(value)
				d.Pct = uint8(pct)
			case "ruf":
				d.P = value
			case "rua":
				d.Rua = value
			case "p":
				d.P = value
			case "sp":
				d.Sp = value
			case "adkim":
				d.Adkim = value
			case "aspf":
				d.Aspf = value
			default:
				// todo: unexpected value, should throw a warning at least
			}
		}
		//TODO: map the dmarc pocliy to the TXT query
	}
	// for _, r := range responses {
	// 	output = append(output, r.String())
	// }
	return
}

func (d DmarcDns) Marshal(kind string) ([]byte, error) {
	switch kind {
	case "json":
		return json.Marshal(d)
	case "yaml":
		return yaml.Marshal(d)
	case "pretty":
		var b bytes.Buffer
		_ = io.Writer(&b)
		dns.PrettyPrint(d, &b, "desc")
		return ioutil.ReadAll(&b)
	case "STIX":
		return nil, fmt.Errorf("STIX has not been implemented yet")
	}
	return nil, fmt.Errorf("Unknown kind: %s", kind)
}
