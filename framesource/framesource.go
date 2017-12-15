package framesource

import (
	"github.com/amarburg/go-frameset/frameset"
	"image"
)

// MovieExtractor is the abstract interface to a quicktime movie.
type FrameSource interface {
	Next() (image.Image, error)
}

func MakeFrameSourceFromPath(  path string ) (FrameSource,error) {

	// Is it a Frameset, a multimov or a movie?

	// Check if it parses as a FrameSet
	set,err := frameset.LoadFrameSet(path)

	if err == nil {
		return MakeFrameSetFrameSource(set)
	}

	return nil,err
}
