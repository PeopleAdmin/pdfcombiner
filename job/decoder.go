// Document structs can either point to a url reference where the pdf is
// stored, or contain a blob of the data itself.  The blob is compressed with
// zlib and then wrapped in Base64.  These functions convert the encoded string
// to a []byte containing the pdf.
package job

import (
	"compress/zlib"
	"encoding/base64"
	"io"
	"io/ioutil"
	"strings"
)

// decodeEmbeddedData takes a base64-encoded string of a gzipped document and
// returns the original source pdf as a byte slice.
func decodeEmbeddedData(encoded string) (decoded []byte, err error) {
	pipeline, err := zlib.NewReader(decoder(encoded))
	defer pipeline.Close()
	if err == nil {
		decoded, err = ioutil.ReadAll(pipeline)
	}
	return
}

func decoder(encoded string) (decoder io.Reader) {
	reader := strings.NewReader(encoded)
	return base64.NewDecoder(base64.StdEncoding, reader)
}
