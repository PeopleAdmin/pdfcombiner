package cpdf

import(
	"io/ioutil"
)

// A standin for a *job.Job.
type bookmarkable interface {
	CombinedBookmarkList() string
	LocalPath() string
	Dir() string
}

// Merge concatenates the given files and adds a header in one pass.
func WriteCombinedBookmarks(job bookmarkable) (err error) {
	infoPath := writeInfoFile(job)
	c := New(job.LocalPath())
	err = c.addBookmarksToFile(infoPath)
	return
}

func (c *Cpdf) addBookmarksToFile(infoPath string) (err error) {
	c.setArgs("-add-bookmarks", infoPath, c.File, "-o", c.File)
	_, err = c.run()
	return
}

func writeInfoFile(job bookmarkable) string {
	infoPath := job.Dir() + "bookmarks.info"
	ioutil.WriteFile(infoPath, []byte(job.CombinedBookmarkList()), 0644)
	return infoPath
}
