package movieset

import (
	"fmt"
	"image"
)

// Thin wrapper around a MovieExtractor which implements Sequental interface
type MovieSequential struct {
	ext      MovieExtractor
	frameNum uint64
}

func MakeMovieSequential(ext MovieExtractor) (*MovieSequential, error) {

	return &MovieSequential{
		ext:      ext,
		frameNum: 1,
	}, nil
}

func (source *MovieSequential) Next() (image.Image, uint64, error) {

	if source.frameNum > source.ext.NumFrames() {
		return nil, source.frameNum, fmt.Errorf("End of file")
	}

	img, err := source.ext.ExtractFrame(source.frameNum)
	source.frameNum++

	return img, source.frameNum - 1, err
}

// func (source *MovieExtractorFrameSource) FrameNum() uint64 {
// 	return source.frameNum
// }
