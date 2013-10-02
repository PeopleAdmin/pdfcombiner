package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

var orig = BookmarkList{
	[]Bookmark{
		Bookmark{0, "1", 1},
		Bookmark{1, "1.a", 2},
		Bookmark{1, "1.b", 3},
	},
}

var other = BookmarkList{
	[]Bookmark{
		Bookmark{0, "1", 2},
		Bookmark{1, "1.a", 3},
		Bookmark{1, "1.b", 4},
	},
}

var empty = BookmarkList{
	[]Bookmark{},
}

func TestBookmarkListInCombinedContext(t *testing.T) {
	assert.Equal(t,
		orig.InCombinedContext("Doc Title", 5),
		BookmarkList{
			[]Bookmark{
				Bookmark{0, "Doc Title", 5},
				Bookmark{1, "1", 5},
				Bookmark{2, "1.a", 6},
				Bookmark{2, "1.b", 7},
			},
		})
}

func TestAlternateBookmarkListInCombinedContext(t *testing.T) {
	assert.Equal(t,
		other.InCombinedContext("Doc Title", 9),
		BookmarkList{
			[]Bookmark{
				Bookmark{0, "Doc Title", 9},
				Bookmark{1, "1", 10},
				Bookmark{2, "1.a", 11},
				Bookmark{2, "1.b", 12},
			},
		})
}

func TestEmptyBookmarkListInCombinedContext(t *testing.T) {
	assert.Equal(t,
		empty.InCombinedContext("Doc Title", 3),
		BookmarkList{
			[]Bookmark{
				Bookmark{0, "Doc Title", 3},
			},
		})
}

func TestBookmarkListStringOutput(t *testing.T) {
	fmt.Println(orig.InCombinedContext("Doc Title", 5).String())
	fmt.Println(orig.InCombinedContext("Doc Title", 5).String())
	assert.Equal(t,
		orig.InCombinedContext("Doc Title", 5).String(),
		`0 "Doc Title" 5
1 "1" 5
2 "1.a" 6
2 "1.b" 7`)
}
