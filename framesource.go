package movieset

import (
	"image"
)

//
type FrameSource interface {
	Next() (image.Image, uint64, error)
	//FrameNum() (uint64)
}

func FrameSourceFromPath(path string) (FrameSource, error) {

	// Is it a Frameset, a multimov or a movie?

	// Check if it parses as a FrameSet
	set, err := LoadFrameSet(path)

	if err == nil {
		return MakeFrameSetFrameSource(set)
	}

	if _, ok := err.(NotAFrameSetError); !ok {
		return nil, err
	}

	ext, err := MovieExtractorFromPath(path)

	if err != nil {
		return nil, err
	}

  return FrameSourceFromMovieExtractor(ext)

}
