package movieset

import (
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"image"
	"path/filepath"
	"time"
)

// MovieExtractor is the abstract interface to a quicktime movie.
type MovieExtractor interface {
	NumFrames() uint64
	Duration() time.Duration
	ExtractFrame(frame uint64) (image.Image, error)
	//ExtractFramePerf(frame uint64) (image.Image, LQTPerformance, error)
}

func OpenMovieExtractor(path string) (MovieExtractor, error) {

	if filepath.Ext(path) == ".mov" {

		file, err := lazyfs.SourceFromPath(path)

		if err != nil {
			return nil, err
		}

		qtInfo, err := lazyquicktime.LoadMovMetadata(file)

		if err != nil {
			return nil, err
		}

		return qtInfo, nil

	} else if filepath.Ext(path) == ".json" {

		mm, err := LoadMultiMov(path)
		if err != nil {
			return nil, err
		}

		return mm, nil

	}

	return nil, fmt.Errorf("Can't make a movie extractor from file %s", path)

}
