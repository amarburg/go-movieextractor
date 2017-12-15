package multimov

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

	lqt *lazyquicktime.LazyQuicktime `json:"-"`
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

// func NewRecordFromFile(file string) (MovRecord, error) {
// 	src, err := lazyfs.OpenLocalFile(file)
// 	if err != nil {
// 		return MovRecord{}, err
// 	}
// 	rec, err := NewRecordFromLazyFs(src)
// 	rec.FullPath = file
// 	return rec, err
// }
//
// func NewRecordFromURL(uri url.URL) (MovRecord, error) {
// 	src, err := lazyfs.OpenHttpSource(uri)
// 	if err != nil {
// 		return MovRecord{}, err
// 	}
// 	rec, err := NewRecordFromLazyFs(src)
// 	rec.FullPath = uri.String()
// 	return rec, err
// }
//
// func NewRecordFromLazyFs(src lazyfs.FileSource) (MovRecord, error) {
// 	qt, err := lqt.LoadMovMetadata(src)
// 	if err != nil {
// 		return MovRecord{}, err
// 	}
// 	rec, err := NewRecord(qt)
// 	rec.ShortName = src.Path()
//
// 	return rec, err
// }
//
// func NewRecord(qt *lqt.LazyQuicktime) (MovRecord, error) {
// 	return MovRecord{}, nil
// }
