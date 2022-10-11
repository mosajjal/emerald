package dkim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/mosajjal/emerald/dns"
	"gopkg.in/yaml.v3"
)

// DKIM is the parsed response for a DKIM query
type DKIM struct {
	D   string `desc:"indicates the domain used with the selector record (s=) to locate the public key"`
	S   string `desc:"indicates the selector record name used with the domain to locate the public key in DNS. The value is a name or number created by the sender"`
	V   string `desc:"is the version of the DKIM record. The value must be DKIM1 and be the first tag in the DNS record. Recommended Optional"`
	K   string `desc:"indicates the key type. The default value is rsa which must be supported by both signers and verifiers."`
	T   string `desc:"indicates the domain is testing DKIM or is enforcing a domain match in the signature header between the i and d tags. Recommended Optional"`
	G   string `desc:"is the granularity of the public key. Optional"`
	H   string `desc:"indicates which hash algorithms are acceptable. Optional"`
	N   string `desc:"is a note field intended for administrators, not end users. Optional"`
	P   string `desc:"indicates the public key of the DKIM record. Required"`
	txt string `desc:"raw TXT response"`
}

// New creates a new DKIM record based on a domain and selector
func New(domain string, selector string) DKIM {
	if strings.HasPrefix(domain, selector) {
		domain = strings.TrimPrefix(domain, selector)
	}
	domain = strings.TrimSuffix(domain, ".")
	// todo: probably more work needed here to be sure domain and prefix are clean
	return DKIM{D: domain, S: selector}
}

// Query function takes a top level domain name (google.com) and
// returns the dkim TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
func (d *DKIM) Query(ctx context.Context, server string) error {
	c, err := dns.NewDnsClient(ctx, server)
	fqdn := d.S + "._domainkey." + d.D
	if err != nil {
		return err
	}
	res, err := c.QueryTXT(ctx, fqdn)
	if len(res) > 0 {
		d.txt = res[0].String()
		// TODO: if the response is a CNAME to another query, that should be handled here
		return d.parseQuery()
	} else {
		return fmt.Errorf("No DKIM record found")
	}
}

// parseQuery function gets the raw TXT query and populates the fields of a DKIM response
// it also checks for errors in the response
func (d *DKIM) parseQuery() error {
	// sample: v=DKIM1; k=rsa; p=KfqqK25Nvy5Gc7t8uGgHW3jJpTxALJqwQIDAQAB
	// ";" looks to be the best way to split the response. All the tests
	// showed that ; is always followed by a space, but we can trim that
	dkimParts := strings.Split(d.txt, ";")
	for _, kv := range dkimParts {
		// trim the whitespace
		kv = strings.TrimSpace(kv)

		// we need to split kv by = to get k and v individually
		tmp := strings.SplitN(kv, "=", 2)
		if len(tmp) != 2 {
			return fmt.Errorf("Wrong TXT response: %s", kv)
		}
		key, value := tmp[0], tmp[1]

		switch key {
		case "v":
			d.V = value
		case "k":
			d.K = value
		case "p":
			d.P = value
		case "h":
			d.H = value
		case "t":
			d.T = value
		case "n":
			d.N = value
		default:
			// todo: unexpected value, should throw a warning at least
		}

	}

	return nil
}

// String function returns the raw TXT response
func (d *DKIM) String() string {
	return d.txt
}

// Marshal provides a way to show a DMARC report in different formats
func (d DKIM) Marshal(kind string) ([]byte, error) {
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
