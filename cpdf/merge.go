package cpdf

import "fmt"

// Acceptable input to Merge(). A standin for a *job.Job.
type mergable interface {
	CombinedTitle() string
	ComponentPaths() []string
	Id() string
	bookmarkable
}

// Merge concatenates the given files, adds a header, and writes combined
// bookmarks in a single pass.
func Merge(job mergable) (err error) {
	err = writeBookmarkInfoFile(job)
	if err != nil {
		return fmt.Errorf("while creating bookmark file: %v", err)
	}
	c := New(job.LocalPath(), job.Id())
	c.addMergeArgs(job.ComponentPaths())
	c.addArgs("AND")
	c.addHeaderArgs(job.CombinedTitle())
	c.addArgs("AND")
	c.addBookmarksArgs(job)
	c.addArgs("-l", "-o", c.File)
	_, err = c.run()
	return
}

func (c *Cpdf) addMergeArgs(paths []string) {
	c.addArgs("-merge")
	for _, path := range paths {
		c.addArgs(path)
	}
}

func (c *Cpdf) addHeaderArgs(title string) {
	headerText := "Page %Page of %EndPage | Created %m-%d-%Y %T | " + title
	c.addArgs("-add-text", headerText, "-top", "15", "-font", "Courier")
}
