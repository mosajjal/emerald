//go:build js

package main

import (
	"encoding/json"
	"github.com/mosajjal/emerald/dmarc"
	"io"
	"log"
	"syscall/js"
)

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("wasmNew", js.FuncOf(newB64))
	<-done
}

func newB64(this js.Value, args []js.Value) interface{} {
	var f io.Reader
	mime, file := dmarc.ParseDataURL(args[0].String())
	if mime == "application/zip" {
		f, _ = dmarc.ParseZipFile(file)
	} else if mime == "application/gzip" {
		f, _ = dmarc.ParseGzipFile(file)
	} else if mime == "text/xml" {
		f = file
	} else {
		log.Fatalln("file format not recognized")
	}

	newReport, _ := dmarc.New(f)
	j, _ := json.MarshalIndent(newReport, "", "  ")

	return string(j)
}
