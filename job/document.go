package job

import (
	"github.com/PeopleAdmin/pdfcombiner/stat"
)

// A Document is a reference to one part of a combined PDF.
// it is is identified by its Key field, which is required.  It can also
// have a Title, used for TOC bookmarks in the final combined document.
// The actual PDF data comes from one of two places:
//  - If the Data field is empty, the Key is treated as a S3 key that is
//    fetched from the enclosing job's bucket.
//  - If the Data field is nonempty, it must contain a zlib-compressed
//    and Base64-encoded string containing the PDF.
type Document struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Data  string `json:"data,omitempty"`
}

// The Key is the only required field
func (doc *Document) isValid() bool {
	return doc.Key != ""
}

// Given a slice of document names, return a slice of Documents.
func docsFromStrings(names []string) (docs []Document) {
	docs = make([]Document, len(names))
	for i, name := range names {
		docs[i] = Document{Key: name}
	}
	return
}

// MarkComplete adds a document to the list of downloaded docs.
// TODO should be a Document.
func (j *Job) MarkComplete(newdoc string, info stat.Stat) {
	j.Downloaded = append(j.Downloaded, newdoc)
	j.PerfStats[newdoc] = info
}

// HasDownloadedDocs determines whether any documents been successfully
// downloaded.
// TODO is it appropriate to use this to determine success in ToJSON()?
func (j *Job) HasDownloadedDocs() bool {
	return len(j.Downloaded) > 0
}
