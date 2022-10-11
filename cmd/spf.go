package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mosajjal/emerald/spf"
)

func init() {
	rootCmd.AddCommand(spfCmd)
	spfCmd.AddCommand(spfQuery)
}

var spfCmd = &cobra.Command{
	Use:   "spf",
	Short: "SPF related commands",
	Long:  `SPF`,
}

var spfQuery = &cobra.Command{
	Use:   "query",
	Short: "Query a Domain's SPF record",
	Long:  "performs a DNS query to a domain's SPF record and parses the output.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), reqTimeout)
		defer cancelFunc()
		mySpf := spf.SpfRecord{QueryDomain: inputDomain}
		if err := mySpf.Query(ctx, dnsServer); err == nil {
			if bytes, err := mySpf.Marshal(outFormat); err == nil {
				os.Stdout.Write(bytes)
			}
		} else {
			fmt.Println(err)
		}
	},
}

func init() {

	spfQuery.PersistentFlags().StringVar(&inputDomain, "domain", "", "input domain. example: google.com")
	spfQuery.MarkPersistentFlagRequired("domain")

}
