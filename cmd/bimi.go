package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
		ctx, cancelFunc := context.WithTimeout(context.Background(), reqTimeout)
		defer cancelFunc()
		if r, img, err := bimi.Query(ctx, inputDomain, dnsServer); err == nil {
			fmt.Println(ConvertImageToANSI(img, 1))
			if out, err := r.Marshal(outFormat); err != nil {
				log.Fatalln(err)
			} else {
				os.Stdout.Write(out)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(bimiCmd)
	bimiCmd.AddCommand(bimiQuery)

	bimiQuery.PersistentFlags().StringVar(&inputDomain, "domain", "", "input domain. example: google.com")
	_ = bimiQuery.MarkPersistentFlagRequired("domain")
}
