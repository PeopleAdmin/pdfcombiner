package job

import (
	"strings"
)

// A Document is a reference to one part of a combined PDF.
// it is is identified by its Key field, which is required.  It can also
// have a Title, used for TOC bookmarks in the final combined document.
// The actual PDF data comes from one of two places:
//  - If the Data field is empty, the Key is treated as a S3 key that is
//    fetched from the enclosing job's bucket.
//  - If the Data field is nonempty, it must contain a zlib-compressed
//    and Base64-encoded string containing the PDF.
// The bookmarks field is filled in by processing and is not serialized.
type Document struct {
	Key       string `json:"key"`
	Title     string `json:"title"`
	Data      string `json:"data,omitempty"`
	PageCount int    `json:"page_count"`
	parent    *Job
	bookmarks BookmarkList
}

func (doc *Document) LocalPath() string {
	return doc.parent.workingDirectory + strings.Replace(doc.Key, "/", "_", -1)
}

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
func (j *Job) MarkComplete(doc *Document) {
	j.Downloaded = append(j.Downloaded, doc)
}

// HasDownloadedDocs determines whether any documents been successfully
// downloaded.
func (j *Job) HasDownloadedDocs() bool {
	return len(j.Downloaded) > 0
}

func (doc *Document) s3Path() string {
	return doc.Key
}
