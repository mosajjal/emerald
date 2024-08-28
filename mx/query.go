package mx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	mkdns "github.com/miekg/dns"
	"github.com/mosajjal/dnsclient"
	"github.com/mosajjal/emerald/dns"
	"gopkg.in/yaml.v3"
)

// TODO:add descriptions for pretty print
type mxRecord struct {
	Priority uint16
	Value    string
	TTL      uint32
}

type MX struct {
	QueryDomain string     `desc:"Queried Domain"`
	Records     []mxRecord `desc:"List of MX records returned"`
}

// Query function takes a top level domain name (google.com) and
// returns the MX TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
func (mx *MX) Query(ctx context.Context, server string) (err error) {
	c, err := dnsclient.New(server, true, "")

	if err != nil {
		return err
	}
	responses, _, err := dns.QueryMX(ctx, c, mx.QueryDomain)
	for _, r := range responses {
		if t, ok := r.(*mkdns.MX); ok {
			tmpMxRecord := mxRecord{}
			tmpMxRecord.Priority = t.Preference
			tmpMxRecord.Value = t.Mx
			tmpMxRecord.TTL = t.Header().Ttl
			mx.Records = append(mx.Records, tmpMxRecord)
		} else {
			return fmt.Errorf("Can't parse MX response")
		}
	}
	return
}

// Marshal provides a way to show a DMARC report in different formats
func (mx MX) Marshal(kind string) ([]byte, error) {
	switch kind {
	case "json":
		return json.Marshal(mx)
	case "yaml":
		return yaml.Marshal(mx)
	case "pretty":
		var b bytes.Buffer
		_ = io.Writer(&b)
		dns.PrettyPrint(mx, &b, "desc")
		return io.ReadAll(&b)
	case "STIX":
		return nil, fmt.Errorf("STIX has not been implemented yet")
	}
	return nil, fmt.Errorf("unknown kind: %s", kind)
}
