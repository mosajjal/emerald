package main

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	outFormat  string
	dnsServer  string
	reqTimeout time.Duration
	rootCmd    = &cobra.Command{
		Use:   "emailtools",
		Short: "Email Parsing tools",
		Long:  `Emailtools is a swiss army knife of dealing with email-related investigations.`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&outFormat, "format", "f", "json", "output format. choices: json, yaml, pretty, stix")
	rootCmd.PersistentFlags().StringVarP(&dnsServer, "dns", "d", "udp://1.1.1.1:53", "DNS server to use. tcp, udp and TLS(tls://9.9.9.9:853) are supported. don't forget the port and URI format!")
	rootCmd.PersistentFlags().DurationVarP(&reqTimeout, "timeout", "t", time.Second*5, "DNS request timeout. example: 100ms, 2s, 50s")
}
