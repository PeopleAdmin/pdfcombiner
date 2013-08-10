package stat

import (
	"time"
)

// A Stat represents statistics about a completed document transfer operation.
type Stat struct {
	Filename string
	Size     int
	DlTime   time.Duration
	Err      error         `json:",omitempty"`
}
