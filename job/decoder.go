package job

import(
	"io/ioutil"
	"io"
	"strings"
	"encoding/base64"
	"compress/zlib"
)

// decodeEmbeddedData takes a base64-encoded string of a gzipped document
// and returns the original source pdf as a byte slice.
func decodeEmbeddedData(encoded string) (decoded []byte, err error) {
	pipeline, err := zlib.NewReader( decoder(encoded) )
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
