package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mosajjal/emerald/dkim"
	"github.com/spf13/cobra"
)

var dkimCmd = &cobra.Command{
	Use:   "dkim",
	Short: "dkim Record",
	Long:  "dkim Record details",
}

var dkimSelector string
var dkimQuery = &cobra.Command{
	Use:   "query",
	Short: "Query a Domain's dkim record",
	Long:  "performs a DNS query to a domain's dkim record and parses the output.",
	Run: func(cmd *cobra.Command, args []string) {
		//todo: write this in dkim's own package and call it from here
		d := dkim.New(inputDomain, dkimSelector)
		ctx, cancelFunc := context.WithTimeout(context.Background(), reqTimeout)
		defer cancelFunc()
		err := d.Query(ctx, dnsServer)
		if err != nil {
			fmt.Println(err)
		} else {
			if out, err := d.Marshal(outFormat); err == nil {
				os.Stdout.Write(out)
			} else {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dkimCmd)
	dkimCmd.AddCommand(dkimQuery)

	dkimQuery.PersistentFlags().StringVar(&inputDomain, "domain", "", "input domain. example: google.com")
	dkimQuery.PersistentFlags().StringVar(&dkimSelector, "selector", "", "DKIM selector")
	dkimQuery.MarkPersistentFlagRequired("selector")
}
