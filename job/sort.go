// The downloaded docs list is generated in parallel, so it's not ordered
// correctly.  These are functions to sort it in order of the original input.
package job

import "sort"

// SortDownloaded ensures that the Downloaded list is in the correct order.
func (j *Job) SortDownloaded() {
	sort.Sort(ByIndex{j.Downloaded})
}

// Given a doc pointer, find its position in DocList.
func (doc *Document) indexOf() int {
	for i, inputDoc := range doc.parent.DocList {
		if inputDoc.Key == doc.Key {
			return i
		}
	}
	panic("Couldn't find " + doc.Key + " in original input!")
}

func (docs Documents) Len() int      { return len(docs) }
func (docs Documents) Swap(i, j int) { docs[i], docs[j] = docs[j], docs[i] }

type ByIndex struct{ Documents }

func (s ByIndex) Less(i, j int) bool { return s.Documents[i].indexOf() < s.Documents[j].indexOf() }
