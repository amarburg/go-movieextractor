package framesource

import (
	"github.com/amarburg/go-frameset/frameset"
	"github.com/amarburg/go-frameset/multimov"
	"image"
)

//
type FrameSource interface {
	Next() (image.Image, uint64, error)
	//FrameNum() (uint64)
}

func MakeFrameSourceFromPath(path string) (FrameSource, error) {

	// Is it a Frameset, a multimov or a movie?

	// Check if it parses as a FrameSet
	set, err := frameset.LoadFrameSet(path)

	if err == nil {
		return MakeFrameSetFrameSource(set)
	}

	if _, ok := err.(frameset.NotAFrameSetError); !ok {
		return nil, err
	}

	ext, err := multimov.MovieExtractorFromPath(path)

	if err != nil {
		return nil, err
	}

  return ext.FrameSource(), nil
	
}
