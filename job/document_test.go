package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocValidation(t *testing.T) {
	var cases = []struct {
		keyName  string
		validity bool
		message  string
	}{
		{"", false, "Empty strings are not valid keys"},
		{"Something", false, "Keys without extensions are not valid"},
		{"pdf", false, "Just because a string has pdf does not mean it is valid"},
		{".pdf", false, "A string without a basename is not a valid key"},
		{"file.pdf", true, "A string with a name and a pdf extension is a valid key"},
	}
	for _, c := range cases {
		doc := &Document{Key: c.keyName}
		assert.Equal(t, doc.isValid(), c.validity, c.message)
	}
}
