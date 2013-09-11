package cpdf

// Acceptable input to Merge(). A standin for a *job.Job.
type mergable interface {
	CombinedTitle() string
	LocalPath() string
	ComponentPaths() []string
}

// Merge concatenates the given files and adds a header in one pass.
func Merge(job mergable) (err error) {
	c := New(job.LocalPath())
	c.addMergeArgs(job.ComponentPaths())
	c.addArgs("AND")
	c.addHeaderArgs(job.CombinedTitle())
	c.addArgs("-o", c.File)
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
