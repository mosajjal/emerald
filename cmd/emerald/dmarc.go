package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/mosajjal/emerald/dmarc"

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
		if err != nil {
			log.Fatalln("filetype cannot be detected. exiting")
		}
		file, err := os.OpenFile(inputFile, os.O_RDONLY, os.FileMode(0o755))
		if err != nil {
			log.Fatalln("cannot open the file, exiting")
		}
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
		if out, err := newReport.Marshal(outFormat); err == nil {
			os.Stdout.Write(out)
		}

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
		ctx, cancelFunc := context.WithTimeout(context.Background(), reqTimeout)
		defer cancelFunc()
		r, err := dmarc.Query(ctx, inputDomain, dnsServer)
		if err != nil {
			log.Fatalln(err)
		}
		out, err := r.Marshal(outFormat)
		if err != nil {
			log.Fatalln(err)
		}
		os.Stdout.Write(out)
	},
}

func init() {
	dmarcParse.PersistentFlags().StringVar(&inputFile, "file", "", "input file")
	dmarcParse.MarkPersistentFlagRequired("file")

	dmarcQuery.PersistentFlags().StringVar(&inputDomain, "domain", "", "input domain. example: google.com")
	dmarcParse.MarkPersistentFlagRequired("domain")

}
