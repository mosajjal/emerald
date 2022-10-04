package dmarc

import (
	"archive/zip"
	"compress/gzip"
	"encoding/base64"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// epoch type is used to get the UNIX timestamp from the XML
// and represent it properly in string or JSON
type epoch time.Time

// ParseDataURL tries to parse JS-provided file uploads in base64 format
// it gets a string blob as an input and outputs the MIME and payload
// example input: data:application/zip;base64,UEsDBAoAAAAIAOhM71TvxiWD5
// example output: "application/zip", an IO reader containing actual bytes
func ParseDataURL(s string) (mime string, content *strings.Reader) {
	r := regexp.MustCompile(`data:(?P<Mime>.*?);base64,(?P<Payload>.*)`)
	match := r.FindStringSubmatch(s)

	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	mime = result["Mime"]
	rString, _ := base64.StdEncoding.DecodeString(result["Payload"])
	content = strings.NewReader(string(rString))
	return
}

// ParseZipFile is a small helper function used to grab
// the first file's handler and return it as an io reader
// essentially streaming unzip function for ZIP format
func ParseZipFile(f io.ReaderAt) (io.Reader, error) {
	gr, err := zip.NewReader(f, 1024)
	if err != nil {
		return nil, err
	}
	return gr.Open(gr.File[0].Name)
}

// ParseGzipFile is a small helper function used to convert
// a GZ io reader into a file io reader, essentially streaming
// gz formatted text back to the caller
func ParseGzipFile(f io.Reader) (io.Reader, error) {
	return gzip.NewReader(f)
}

func (e *epoch) UnmarshalText(data []byte) error {
	if i, err := strconv.Atoi(string(data)); err == nil {
		*(*time.Time)(e) = time.Unix(int64(i), 0)
	} else {
		return err
	}
	return nil
}
func (e *epoch) UnmarshalBinary(data []byte) error {
	if i, err := strconv.Atoi(string(data)); err == nil {
		*(*time.Time)(e) = time.Unix(int64(i), 0)
	} else {
		return err
	}
	return nil
}

func (e epoch) MarshalBinary() ([]byte, error) { return (time.Time)(e).MarshalBinary() }
func (e epoch) MarshalText() ([]byte, error)   { return time.Time(e).MarshalText() }
