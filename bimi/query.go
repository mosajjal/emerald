package bimi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	mkdns "github.com/miekg/dns"
	"gopkg.in/yaml.v3"

	"github.com/mosajjal/emerald/dns"

	"github.com/mosajjal/dnsclient"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type BimiDns struct {
	V string `desc:"BIMI version"`
	L string `desc:"SVG URL"`
	A string `desc:"BIMI VMC PEM URL"`
}

func getImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	icon, err := oksvg.ReadIconStream(resp.Body)
	if err != nil {
		return nil, err
	}
	icon.SetTarget(0, 0, 128, 128)
	rgba := image.NewRGBA(image.Rect(0, 0, 128, 128))
	icon.Draw(rasterx.NewDasher(128, 128, rasterx.NewScannerGV(128, 128, rgba, rgba.Bounds())), 1)

	return rgba, nil
}

// Query function takes a top level domain name (google.com) and
// returns the bimi TXT record associated with it. it uses the
// system's resolver if server is provided as 0.0.0.0 otherwise
// it'll explicity query from the requested server.
func Query(ctx context.Context, domain string, server string) (d BimiDns, img image.Image, err error) {
	if !strings.HasPrefix(domain, "_bimi.") {
		domain = "default._bimi." + domain
	}
	domain = strings.TrimSuffix(domain, ".")
	// c, err := dns.NewDnsClient(ctx, server)
	c, err := dnsclient.New(server, true)
	if err != nil {
		return
	}
	responses, _, err := dns.QueryTXT(ctx, c, domain)
	if err != nil {
		return
	}
	// parse TXT response and grab the URL if exists (l=https://bimigroup.org/bimi-sq.svg;)
	//TODO: in next line we only take the first output into account, there might be a chance
	// there are multiple outputs
	if t, ok := responses[0].(*mkdns.TXT); ok {
		bimiParts := strings.Split(t.Txt[0], ";")
		for _, kv := range bimiParts {
			// trim the whitespace
			kv = strings.TrimSpace(kv)
			tmp := strings.SplitN(kv, "=", 2)
			if len(tmp) != 2 {
				continue
			}
			key, value := tmp[0], tmp[1]
			switch key {
			case "v":
				d.V = value
			case "l":
				d.L = value
			case "a":
				d.A = value
			}
		}
	}
	img, err = getImageFromURL(d.L)
	return
}

func (d BimiDns) Marshal(kind string) ([]byte, error) {
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
