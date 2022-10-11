package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mosajjal/emerald/mx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mxCmd)
	mxCmd.AddCommand(mxQuery)
}

var mxCmd = &cobra.Command{
	Use:   "mx",
	Short: "MX related commands",
	Long:  `MX`,
}

var mxQuery = &cobra.Command{
	Use:   "query",
	Short: "Query a Domain's MX record",
	Long:  "performs a DNS query to a domain's MX record and parses the output.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), reqTimeout)
		defer cancelFunc()
		myMx := mx.MX{QueryDomain: inputDomain}
		if err := myMx.Query(ctx, dnsServer); err == nil {
			if bytes, err := myMx.Marshal(outFormat); err == nil {
				os.Stdout.Write(bytes)
			}
		} else {
			fmt.Println(err.Error())
		}
	},
}

func init() {

	mxQuery.PersistentFlags().StringVar(&inputDomain, "domain", "", "input domain. example: google.com")
	mxQuery.MarkPersistentFlagRequired("domain")

}
