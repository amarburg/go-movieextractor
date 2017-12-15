package multimov

import (
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"image"
	"log"
	"os"
	"time"
)

// MultiMov.NumFrames is defined in the main MultiMov API

// Duration calculates the total continuous duration within a MultiMov.
// The sum of all movie durations, does account for time gaps
func (mm MultiMov) Duration() time.Duration {
	var out time.Duration
	for _, h := range mm.Sequence {
		out += mm.Movies[h.Hash].Duration
	}

	return out
}

// ExtractFrame extracts the specified frame from a MultiMov
func (mm MultiMov) ExtractFrame(frame uint64) (image.Image, error) {
	hash, offset, err := mm.Offset(frame)

	if err != nil {
		return image.NewGray(image.Rect(0, 0, 0, 0)), err
	}

	mov, has := mm.Movies[hash]

	if !has {
		return image.NewGray(image.Rect(0, 0, 0, 0)), fmt.Errorf("Error looking up movie %x in table", hash)
	}

	if mov.lqt == nil {
		movFile := mm.MovPath(hash)

		if _, err := os.Stat(movFile); os.IsNotExist(err) {
			return nil, err
		}

		log.Printf("Opening movie file: %s", movFile)

		fs, err := lazyfs.OpenLocalFile(movFile)
		if err != nil {
			return nil, err
		}

		lqt, err := lazyquicktime.LoadMovMetadata(fs)
		if err != nil {
			return nil, err
		}

		mov.lqt = lqt
		mm.Movies[hash] = mov
	}

	return mov.lqt.ExtractFrame(frame - offset)
}

// ExtractFrame extracts the specified frame from a MultiMov
func (mm MultiMov) ExtractFramePerf(frame uint64) (image.Image, lazyquicktime.LQTPerformance, error) {
	hash, offset, err := mm.Offset(frame)

	if err != nil {
		return image.NewGray(image.Rect(0, 0, 0, 0)), lazyquicktime.LQTPerformance{}, err
	}

	mov, has := mm.Movies[hash]

	if !has {
		return image.NewGray(image.Rect(0, 0, 0, 0)), lazyquicktime.LQTPerformance{}, fmt.Errorf("Error looking up movie %x in table", hash)
	}

	if mov.lqt == nil {
		movFile := mm.MovPath(hash)

		if _, err := os.Stat(movFile); os.IsNotExist(err) {
			return nil, lazyquicktime.LQTPerformance{}, err
		}

		log.Printf("Opening movie file: %s", movFile)

		fs, err := lazyfs.OpenLocalFile(movFile)
		if err != nil {
			return nil, lazyquicktime.LQTPerformance{}, err
		}

		lqt, err := lazyquicktime.LoadMovMetadata(fs)
		if err != nil {
			return nil, lazyquicktime.LQTPerformance{}, err
		}

		mov.lqt = lqt
		mm.Movies[hash] = mov
	}

	return mov.lqt.ExtractFramePerf(frame - offset)
}
