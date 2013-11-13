package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var docs = []Document{
	{
		Title:     "First",
		PageCount: 5,
		Bookmarks: BookmarkList{
			[]Bookmark{
				{0, "1", 1},
				{1, "1.a", 2},
				{1, "1.b", 3},
			},
		},
	},
	{
		Title:     "Second",
		PageCount: 3,
		Bookmarks: BookmarkList{[]Bookmark{}},
	},
	{
		Title:     "Third",
		PageCount: 8,
		Bookmarks: BookmarkList{
			[]Bookmark{
				{0, "1", 1},
				{1, "1.a", 2},
				{2, "1.a.i", 3},
				{0, "2", 4},
			},
		},
	},
}

var job = Job{DocList: docs}

func TestJobCombinedBookmarkList(t *testing.T) {
	assert.Equal(t,
		job.CombinedBookmarkList(), `0 "First" 1
1 "1" 1
2 "1.a" 2
2 "1.b" 3
0 "Second" 6
0 "Third" 9
1 "1" 9
2 "1.a" 10
3 "1.a.i" 11
1 "2" 12`)
}
