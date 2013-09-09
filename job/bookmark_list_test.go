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

func TestBookmarkListInCombinedContext(t *testing.T) {
	assert.Equal(t,
		orig.InCombinedContext("Doc Title", 5),
		BookmarkList{
			[]Bookmark{
				Bookmark{0, "Doc Title", 5},
				Bookmark{1, "1", 6},
				Bookmark{2, "1.a", 7},
				Bookmark{2, "1.b", 8},
			},
		})
}

func TestBookmarkListStringOutput(t *testing.T) {
	fmt.Println(orig.InCombinedContext("Doc Title", 5).String())
	fmt.Println(orig.InCombinedContext("Doc Title", 5).String())
	assert.Equal(t,
		orig.InCombinedContext("Doc Title", 5).String(),
		`0 "Doc Title" 5
1 "1" 6
2 "1.a" 7
2 "1.b" 8`)
}
