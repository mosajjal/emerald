package dmarc

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"

	"github.com/mosajjal/emerald/dns"
	yaml "gopkg.in/yaml.v3"
)

type dateRange struct {
	Begin epoch `xml:"begin"`
	End   epoch `xml:"end"`
}

type dmarcMetadata struct {
	OrgName   string    `xml:"org_name"`
	Email     string    `xml:"email"`
	ReportID  string    `xml:"report_id"`
	DateRange dateRange `xml:"date_range"`
}

// v	DMARC1	DMARC protocol version.
// p	reject	Apply this policy to email that fails the DMARC check. This policy be set to 'none', 'quarantine', or 'reject'. 'none' is used to collect the DMARC report and gain insight into the current emailflows and their status.
// pct	100	The percentage tag instructs ISPs to only apply the DMARC policy to a percentage of failing email's. 'pct = 50' will tell receivers to only apply the 'p = ' policy 50% of the time against email's that fail the DMARC check. NOTE: this will not work for the 'none' policy, but only for 'quarantine' or 'reject' policies.
// rua	mailto:hi@n0p.me	A list of URIs for ISPs to send XML feedback to. NOTE: this is not a list of email addresses. DMARC requires a list of URIs of the form 'mailto:test@example.com'.
// ruf		A list of URIs for ISPs to send forensic reports to. NOTE: this is not a list of email addresses. DMARC requires a list of URIs of the form 'mailto:test@example.org'.
// rf	afrf	The reporting format for forensic reports. This can be either 'afrf' or 'iodef'.
// adkim	r	Specifies the 'Alignment Mode' for DKIM signatures, this can be either 'r' (Relaxed) or 's' (Strict). In Relaxed mode also authenticated DKIM signing domains (d=) that share a Organizational Domain with an emails From domain will pass the DMARC check. In Strict mode an exact match is required.
// aspf	r	Specifies the 'Alignment Mode' for SPF, this can be either 'r' (Relaxed) or 's' (Strict). In Relaxed mode also authenticated SPF domains that share a Organizational Domain with an emails From domain will pass the DMARC check. In Strict mode an exact match is required.
// sp	p=value	This policy should be applied to email from a sub-domain of this domain that fail the DMARC check. Using this tag domain owners can publish a 'wildcard' policy for all subdomains.
// fo	0	Forensic options. Allowed values: '0' to generate reports if both DKIM and SPF fail, '1' to generate reports if either DKIM or SPF fails to produce a DMARC pass result, 'd' to generate report if DKIM has failed or 's' if SPF failed.
// ri	86400	The reporting interval for how often you'd like to receive aggregate XML reports. This is a preference and ISPs could (and most likely will) send the report on different intervals (normally this will be daily).
type dmarcPolicy struct {
	Domain string `xml:"domain"`
	Adkim  string `xml:"adkim"`
	Aspf   string `xml:"aspf"`
	P      string `xml:"p"`
	Sp     string `xml:"sp"`
	Pct    uint8  `xml:"pct"`
	Fo     uint8  `xml:"fo"`
}

type recordPolicy struct {
	Disposition string `xml:"disposition"`
	Fail        string `xml:"dkim"`
	Spf         string `xml:"spf"`
}

type recordRow struct {
	SourceIP   net.IP       `xml:"source_ip"`
	Count      uint64       `xml:"count"`
	PolicyEval recordPolicy `xml:"policy_evaluated"`
}

type recordID struct {
	EnvelopeTo   string `xml:"envelope_to"`
	EnvelopeFrom string `xml:"envelope_from"`
	HeaderFrom   string `xml:"header_from"`
}

type authResSpf struct {
	Domain string `xml:"domain"`
	Scope  string `xml:"scope"`
	Result string `xml:"result"`
}
type authResDkim struct {
	Domain      string `xml:"domain"`
	Result      string `xml:"result"`
	HumanResult string `xml:"human_result"`
}

type recordAuthRes struct {
	Spf  authResSpf  `xml:"spf"`
	Dkim authResDkim `xml:"dkim"`
}

type dmarcRecord struct {
	Row         recordRow     `xml:"row"`
	Identifiers recordID      `xml:"identifiers"`
	AuthResults recordAuthRes `xml:"auth_results"`
}

// AfrfReport is the outline of the XML
type AfrfReport struct {
	//XMLName  xml.Name      `xml:"feedback"`
	Version  string        `xml:"version"`
	Metadata dmarcMetadata `xml:"report_metadata"`
	Policy   dmarcPolicy   `xml:"policy_published"`
	Records  []dmarcRecord `xml:"record"`
}

// New will be exported to WASM
func New(f io.Reader) (AfrfReport, error) {
	var res []byte
	buf := make([]byte, 2048)
	for {
		n, err := f.Read(buf)
		res = append(res, buf[0:n]...)
		if err != nil {
			break
		}
	}
	newReport := new(AfrfReport)
	err := xml.Unmarshal(res, newReport)
	return *newReport, err
}

// Marshal provides a way to show a DMARC report in different formats
func (afrf AfrfReport) Marshal(kind string) ([]byte, error) {
	switch kind {
	case "json":
		return json.Marshal(afrf)
	case "yaml":
		return yaml.Marshal(afrf)
	case "pretty":
		var b bytes.Buffer
		_ = io.Writer(&b)
		dns.PrettyPrint(afrf, &b, "desc")
		return io.ReadAll(&b)
	case "STIX":
		return nil, fmt.Errorf("STIX has not been implemented yet")
	}
	return nil, fmt.Errorf("unknown kind: %s", kind)
}
