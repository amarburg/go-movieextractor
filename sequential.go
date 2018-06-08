package movieset

import (
	"image"
)

// Interface
type Sequential interface {
	Next() (image.Image, uint64, error)
	//FrameNum() (uint64)
}

func OpenSequential(path string) (Sequential, error) {

	// Is it a Frameset, a multimov or a movie?

	// Check if it parses as a FrameSet
	set, err := LoadFrameSet(path)
	if err == nil {
		return MakeFrameSetSequential(set)
	}

	if _, ok := err.(NotAFrameSetError); !ok {
		return nil, err
	}

	ext, err := OpenMovieExtractor(path)

	if err != nil {
		return nil, err
	}

	return MakeMovieSequential(ext)

}
