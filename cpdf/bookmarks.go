package cpdf

import (
	"io/ioutil"
)

// A standin for a *job.Job.
type bookmarkable interface {
	CombinedBookmarkList() string
	LocalPath() string
	Dir() string
}

func (c *Cpdf) addBookmarksArgs(job bookmarkable) {
	c.addArgs("-add-bookmarks", infoPath(job))
}

func writeBookmarkInfoFile(job bookmarkable) error {
	return ioutil.WriteFile(infoPath(job), []byte(job.CombinedBookmarkList()), 0644)
}

func infoPath(job bookmarkable) string {
	return job.Dir() + "bookmarks.info"
}
