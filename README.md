
# CLI Email toolkit

## Sections

### DMARC
- Query DMARC of a domain and display the parsed output
- Parse DMARC XML report from .xml, .zip and .gz formats

### DKIM
- Query DKIM of a domain and display the parsed output

### SPF
- Query SPF of a domain and display the parsed output

### MX
- Query MX of a domain and display the parsed output

### MTA-STS 
- Query MTA-STS of a domain and display the parsed output

### BIMI
- Query BIMI of a domain and display the parsed output
- Verify VMC
- Verify the SVG format

### EMAIL
- convert MSG to EML -> too hard for v1, will consider for v2
- parse EML file and return:
    - IoCs
    - real sender
    - highlighted output
    - confidence level of spam
    - Graph of servers?
    - MHA link for more info
- verify EML's DKIM

## Output Formats
- JSON
- YAML 
- Pretty colorful (default)
- Pretty raw
- STIXX




Ideas:

generate the appropiate DKIM query to test against an EML file for DKIM 

