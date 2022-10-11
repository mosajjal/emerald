
# Email security toolkit


`emerald` is a WIP project aimed at consolidating different aspects of email security in one place. It aims to be a similar tool to MX toolbox but from your own network.

The following are implemented

### DMARC
- Query DMARC of a domain and display the parsed output
- Parse DMARC XML report from .xml, .zip and .gz formats

### DKIM
- Query DKIM of a domain and display the parsed output
```bash
$ emerald  dkim query --domain n0p.me --selector 1
┌───────┬────────────────────────────────────────┬────────────────────────────────────────┐
│ PARAM │ DESC                                   │ VALUE                                  │
├───────┼────────────────────────────────────────┼────────────────────────────────────────┤
│ D     │ indicates the domain used with the     │ n0p.me                                 │
│       │ selector record (s=) to locate the     │                                        │
│       │ public key                             │                                        │
├───────┼────────────────────────────────────────┼────────────────────────────────────────┤
│ S     │ indicates the selector record name     │ 1                                      │
│       │ used with the domain to locate the     │                                        │
│       │ public key in DNS. The value is a name │                                        │
│       │ or number created by the sender        │                                        │
├───────┼────────────────────────────────────────┤                                        │
│ V     │ is the version of the DKIM record. The │                                        │
│       │ value must be DKIM1 and be the first   │                                        │
│       │ tag in the DNS record. Recommended     │                                        │
│       │ Optional                               │                                        │
├───────┼────────────────────────────────────────┼────────────────────────────────────────┤
│ K     │ indicates the key type. The default    │ rsa                                    │
│       │ value is rsa which must be supported   │                                        │
│       │ by both signers and verifiers.         │                                        │
├───────┼────────────────────────────────────────┤                                        │
│ T     │ indicates the domain is testing DKIM   │                                        │
│       │ or is enforcing a domain match in the  │                                        │
│       │ signature header between the i and d   │                                        │
│       │ tags. Recommended Optional             │                                        │
├───────┼────────────────────────────────────────┤                                        │
│ G     │ is the granularity of the public key.  │                                        │
│       │ Optional                               │                                        │
├───────┼────────────────────────────────────────┤                                        │
│ H     │ indicates which hash algorithms are    │                                        │
│       │ acceptable. Optional                   │                                        │
├───────┼────────────────────────────────────────┤                                        │
│ N     │ is a note field intended for           │                                        │
│       │ administrators, not end users.         │                                        │
│       │ Optional                               │                                        │
├───────┼────────────────────────────────────────┼────────────────────────────────────────┤
│ P     │ indicates the public key of the DKIM   │ MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQ │
│       │ record. Required                       │ C2VfTfGCVlBkPfm5VeKwjQtoQvRvmUyaR4UWrE │
│       │                                        │ WYdBwUdkrszuqj4Vkp6w6clTyxTlw+RdfYtlab │
│       │                                        │ f3CKWldHeXhjcE+C2lwspdJR0YlcWlOHSXT1Eo │
│       │                                        │ 1c8QYHfzhruKK53qg/cScsWM+KfqqK25Nvy5Gc │
│       │                                        │ 7t8uGgHW3jJpTxALJqwQIDAQAB"            │
└───────┴────────────────────────────────────────┴────────────────────────────────────────┘

```

### SPF
- Query SPF of a domain and display the parsed output
```bash
$ emerald  spf query --domain google.com -d tcp://9.9.9.9:53
┌─────────────┬────────────────────────────────────────┬─────────────────┐
│ PARAM       │ DESC                                   │ VALUE           │
├─────────────┼────────────────────────────────────────┼─────────────────┤
│ QueryDomain │ Domain                                 │ google.com      │
├─────────────┼────────────────────────────────────────┼─────────────────┤
│ Version     │ SPF record version                     │ spf1            │
├─────────────┼────────────────────────────────────────┼─────────────────┤
│ Includes    │ List of allowed third party domains    │ _spf.google.com │
│             │ that can send email on you behalf      │                 │
├─────────────┼────────────────────────────────────────┼─────────────────┤
│ Tag         │ Tag (policy) applied to SPF. -all:     │ ~all            │
│             │ fail, ~all: softfail, and +all: allow  │                 │
└─────────────┴────────────────────────────────────────┴─────────────────┘

```

### MX
- Query MX of a domain and display the parsed output
```bash
$ emerald mx query --domain n0p.me -d udp://9.9.9.9:53

┌─────────────┬─────────────────────────────┬──────────┬───────────────┐
│ PARAM       │ DESC                        │ VALUE    │               │
├─────────────┼─────────────────────────────┼──────────┼───────────────┤
│ QueryDomain │ Queried Domain              │ n0p.me   │               │
├─────────────┼─────────────────────────────┼──────────┼───────────────┤
│ Records     │ List of MX records returned │ Priority │ 10            │
│             │                             ├──────────┼───────────────┤
│             │                             │ Value    │ mx.zoho.com.  │
│             │                             ├──────────┼───────────────┤
│             │                             │ TTL      │ 16            │
│             │                             ├──────────┼───────────────┤
│             │                             │ Priority │ 20            │
│             │                             ├──────────┼───────────────┤
│             │                             │ Value    │ mx2.zoho.com. │
│             │                             ├──────────┼───────────────┤
│             │                             │ TTL      │ 16            │
└─────────────┴─────────────────────────────┴──────────┴───────────────┘
```


### MTA-STS 
- Query MTA-STS of a domain and display the parsed output (WIP)

### BIMI
- Query BIMI of a domain and display the parsed output
- Verify VMC (WIP)
- Verify the SVG format (WIP)
- Display the SVG in the terminal 
```bash
$ emerald bimi query --domain c-date.de

[SVG is shown here in terminal]

┌───────┬──────────────────┬────────────────────────────────────────┐
│ PARAM │ DESC             │ VALUE                                  │
├───────┼──────────────────┼────────────────────────────────────────┤
│ V     │ BIMI version     │ BIMI1                                  │
├───────┼──────────────────┼────────────────────────────────────────┤
│ L     │ SVG URL          │ https://bimi.entrust.net/c-date.de/log │
│       │                  │ o.svg                                  │
├───────┼──────────────────┼────────────────────────────────────────┤
│ A     │ BIMI VMC PEM URL │ https://bimi.entrust.net/c-date.de/cer │
│       │                  │ tchain.pem                             │
└───────┴──────────────────┴────────────────────────────────────────┘

```
## Output Formats
- JSON
- YAML 
- Pretty colorful (default)
- STIX (WIP)

