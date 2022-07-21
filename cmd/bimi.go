package cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/mosajjal/emerald/bimi"
	"github.com/spf13/cobra"
)

var bimiCmd = &cobra.Command{
	Use:   "bimi",
	Short: "BIMI Record",
	Long:  "BIMI Record details",
}

var bimiQuery = &cobra.Command{
	Use:   "query",
	Short: "Query a Domain's bimi record",
	Long:  "performs a DNS query to a domain's bimi record and parses the output.",
	Run: func(cmd *cobra.Command, args []string) {
		//todo: write this in bimi's own package and call it from here
		r, _ := bimi.Query(context.Background(), inputDomain, net.IPv4zero)
		fmt.Println(r)
	},
}

func init() {
	rootCmd.AddCommand(bimiCmd)
	bimiCmd.AddCommand(bimiQuery)

	bimiCmd.PersistentFlags().StringVar(&inputDomain, "domain", "", "input domain. example: google.com")
}
