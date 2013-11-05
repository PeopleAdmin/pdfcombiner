package job

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/cpdf"
	"regexp"
	"strconv"
)

var bookmarkPattern = regexp.MustCompile(`^(\d+) "(.*)" (\d+)`)

type Bookmark struct {
	Depth int
	Name  string
	Page  int
}

func (doc *Document) GetMetadata(cmd cpdf.InfoCmd) (err error) {
	err = cmd.Validate()
	if err != nil {
		doc.parent.AddError(fmt.Errorf("Validating %v: %v", doc.Key, err))
		return
	}
	doc.PageCount, err = cmd.PageCount()
	if err != nil {
		doc.parent.AddError(fmt.Errorf("Counting Pages for %v: %v", doc.Key, err))
		return
	}
	doc.Bookmarks, err = ExtractBookmarks(cmd)
	if err != nil {
		doc.parent.AddError(fmt.Errorf("Extracting Bookmarks for %v: %v", doc.Key, err))
		return
	}
	return
}

func (b *Bookmark) String() string {
	return fmt.Sprintf("%d %s %d", b.Depth, strconv.Quote(b.Name), b.Page)
}

func ExtractBookmarks(cmd cpdf.InfoCmd) (bookmarks BookmarkList, err error) {
	bookmarks = BookmarkList{}
	raw, err := cmd.ListBookmarks()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return
		}
		addBookmark(scanner.Text(), &bookmarks)
	}
	return
}

func addBookmark(cpdfOutputLine string, bookmarks *BookmarkList) {
	bm, err := newBookmark(cpdfOutputLine)
	if err != nil {
		return
	}
	bookmarks.add(bm)
}

func newBookmark(bmLine string) (bm Bookmark, err error) {
	err = fmt.Errorf("Invalid Bookmark Format: " + bmLine)
	matches := bookmarkPattern.FindStringSubmatch(bmLine)
	if len(matches) != 4 {
		return
	}
	depth, e := strconv.Atoi(matches[1])
	if e != nil {
		return
	}
	page, e := strconv.Atoi(matches[3])
	if e != nil {
		return
	}
	err = nil
	bm = Bookmark{
		Depth: depth,
		Name:  matches[2],
		Page:  page,
	}
	return
}
