package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/mosajjal/emerald/dmarc"
	"github.com/mosajjal/emerald/dns"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dmarcCmd)
	dmarcCmd.AddCommand(dmarcParse)
	dmarcCmd.AddCommand(dmarcQuery)
}

var dmarcCmd = &cobra.Command{
	Use:   "dmarc",
	Short: "DMARC related commands",
	Long:  `DMARC command set has tools to query, parse and investingate DMARC related issues`,
}

var inputFile string
var dmarcParse = &cobra.Command{
	Use:   "parse",
	Short: "Parse DMARC Reports",
	Long:  `used to parse any DMARC XML reports in .xml, .zip or .gz formats`,
	Run: func(cmd *cobra.Command, args []string) {
		mtype, err := mimetype.DetectFile(inputFile)
		file, err := os.OpenFile(inputFile, os.O_RDONLY, os.FileMode(0o755))
		defer file.Close()
		if err != nil {
			log.Fatalln(err)
		}
		var f io.Reader

		if mtype.Is("application/zip") {
			f, _ = dmarc.ParseZipFile(file)
		} else if mtype.Is("application/gzip") {
			f, _ = dmarc.ParseGzipFile(file)
		} else if mtype.Is("text/xml") {
			f = file
		} else {
			log.Fatalln("file format not recognized")
		}

		newReport, _ := dmarc.New(f)
		dns.PrettyPrint(newReport, os.Stdout)

		//j, _ := json.MarshalIndent(newReport, "", "  ")
		//fmt.Printf("%s", j)
	},
}

var inputDomain string
var dmarcQuery = &cobra.Command{
	Use:   "query",
	Short: "Query a Domain's DMARC record",
	Long:  "performs a DNS query to a domain's DMARC record and parses the output.",
	Run: func(cmd *cobra.Command, args []string) {
		//todo: write this in DMARC's own package and call it from here
		r, e := dmarc.Query(context.Background(), inputDomain, net.IPv4zero)
		fmt.Println(r, e)
	},
}

func init() {
	dmarcParse.PersistentFlags().StringVar(&inputFile, "file", "", "input file")
	dmarcParse.MarkPersistentFlagRequired("file")

	dmarcQuery.PersistentFlags().StringVar(&inputDomain, "domain", "", "input domain. example: google.com")
	dmarcParse.MarkPersistentFlagRequired("domain")

}
