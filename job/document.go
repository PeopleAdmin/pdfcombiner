package job

import (
	"github.com/PeopleAdmin/pdfcombiner/stat"
)

// A Document has a (file) name and a human readable title, possibly
// used for watermarking prior to combination.
type Document struct {
	Key   string `json:"key,omitempty"`
	Title string `json:"title"`
	Data  string `json:"data,omitempty"`
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
