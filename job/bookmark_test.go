package job

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var validCpdfOutput = `0 "Bookmark One" 1
0 "Bookmark Two" 2
1 "Nested under Two" 3
0 "Bad chars in title: \" 123" 4
`
var expectedParseOutput = []Bookmark{
	{0, "Bookmark One", 1},
	{0, "Bookmark Two", 2},
	{1, "Nested under Two", 3},
	{0, `Bad chars in title: \" 123`, 4},
}

// Set up serveral mock varities of the cpdf.InfoCmd type.
type validOutputter struct{}
type invalidOutputter struct{}
type erroringOutputter struct{}

func (o *invalidOutputter) PageCount() (int, error)  { return 0, nil }
func (o *validOutputter) PageCount() (int, error)    { return 0, nil }
func (o *erroringOutputter) PageCount() (int, error) { return 0, nil }

func (o *validOutputter) ListBookmarks() (out []byte, err error) {
	out = []byte(validCpdfOutput)
	return
}
func (o *invalidOutputter) ListBookmarks() (out []byte, err error) {
	out = []byte("some invalid output")
	return
}
func (o *erroringOutputter) ListBookmarks() (out []byte, err error) {
	err = fmt.Errorf("ListBookmarks() threw an error!")
	return
}

func (o *invalidOutputter) Validate() error  { return errors.New("") }
func (o *validOutputter) Validate() error    { return errors.New("") }
func (o *erroringOutputter) Validate() error { return errors.New("") }

func TestParsesInvalidBookmarks(t *testing.T) {
	bookmarks, _ := ExtractBookmarks(&invalidOutputter{})
	assert.Equal(t, bookmarks.len(), 0)
}

func TestHandlesPdfErrors(t *testing.T) {
	bookmarks, err := ExtractBookmarks(&erroringOutputter{})
	assert.Error(t, err)
	assert.Equal(t, bookmarks.len(), 0)
}

func TestParsesValidBookmarks(t *testing.T) {
	bookmarks, err := ExtractBookmarks(&validOutputter{})
	assert.NoError(t, err)
	assert.Equal(t, len(expectedParseOutput), bookmarks.len())
	assert.Equal(t, bookmarks.list, expectedParseOutput)
}

func TestBookmarkStringification(t *testing.T) {
	var cases = []struct {
		in      Bookmark
		out     string
		message string
	}{
		{Bookmark{0, `Blah`, 1}, `0 "Blah" 1`, "Normal bookmark output"},
		{Bookmark{0, `Bl"ah`, 1}, `0 "Bl\"ah" 1`, "Quotes contained in bookmark titles are properly escaped"},
	}
	for _, c := range cases {
		assert.Equal(t, c.in.String(), c.out, c.message)
	}
}
