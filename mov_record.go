package movieset

import (
	"github.com/amarburg/go-lazyquicktime"
	"path"
	"time"
)

// MovRecord stores filepath and meta-information about a single Mov within a
// Multimov
type MovRecord struct {
	ShortName string        `json:",omitempty"`
	Relapath  string        `json:",omitempty"`
	NumFrames uint64        `json:",omitempty"`
	StartTime time.Time     `json:",omitempty"`
	Duration  time.Duration `json:",omitempty"`

	lqt *lazyquicktime.LazyQuicktime
}

// MovRecordFromLqt creates a MovRecord given a LazyQuicktime
func MovRecordFromLqt(lqt *lazyquicktime.LazyQuicktime) MovRecord {
	return MovRecord{
		ShortName: path.Base(lqt.Source.Path()),
		Relapath:  lqt.Source.Path(),
		NumFrames: lqt.NumFrames(),
		Duration:  time.Duration(lqt.Duration() * time.Second),
		lqt:       lqt,
	}
}
