package job

import(
	"strings"
)

type BookmarkList struct {
	list []Bookmark
}

// InCombinedContext returns a new BookmarkList representing the document's
// bookmarks as they would appear in the context of a combined document.  That
// is, a new level-0 bookmark is added with all existing bookmarks as children,
// and all pages are offset by the document's start position.
func (l *BookmarkList) InCombinedContext(title string, offset int) (newBm BookmarkList) {
	newBm = BookmarkList{
		[]Bookmark{
			Bookmark{0, title, offset},
		},
	}
	newBm.concat(l.oneDeeper().offsetBy(offset))
	return
}

func (l BookmarkList) String() string {
	bookmarkStrings := make([]string, l.len())
	for i, bm := range l.list {
		bookmarkStrings[i] = bm.String()
	}
	return strings.Join(bookmarkStrings, "\n")
}

// A new *BookmarkList with each Bookmark's Depth incremented.
func (l BookmarkList) oneDeeper() *BookmarkList {
	newList := BookmarkList{make([]Bookmark, l.len())}
	copy(newList.list, l.list)
	for i, _ := range l.list {
		newList.list[i].Depth = l.list[i].Depth + 1
	}
	return &newList
}

// A new *BookmarkList with each Bookmark's Page offset by the given amount.
func (l BookmarkList) offsetBy(offset int) *BookmarkList {
	newList := BookmarkList{make([]Bookmark, l.len())}
	copy(newList.list, l.list)
	for i, _ := range l.list {
		newList.list[i].Page = l.list[i].Page + offset
	}
	return &newList
}

func (l *BookmarkList) concat(other *BookmarkList) {
	l.list = append(l.list, other.list...)
}

func (l *BookmarkList) add(bm Bookmark) {
	l.list = append(l.list, bm)
}

func (l *BookmarkList) len() int {
	return len(l.list)
}
