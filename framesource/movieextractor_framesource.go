package framesource

import (
	"fmt"
	"github.com/amarburg/go-lazyquicktime"
	"image"
)

// Thin wrapper around a MovieExtractor which implements FrameSource
type MovieExtractorFrameSource struct {
	lazyquicktime.MovieExtractor
	frameNum uint64
}


func (ext lazyquicktime.MovieExtractor) FrameSource() (*MovieExtractorFrameSource, error) {
	return &MovieExtractorFrameSource{
		MovieExtractor: ext,
		frameNum:       1,
	}, nil
}

// func MakeMovieExtractorFrameSource(ext lazyquicktime.MovieExtractor) (*MovieExtractorFrameSource, error) {
//
// 	return &MovieExtractorFrameSource{
// 		MovieExtractor: ext,
// 		frameNum:       1,
// 	}, nil
// }

func (source *MovieExtractorFrameSource) Next() (image.Image, uint64, error) {

	if source.frameNum > source.MovieExtractor.NumFrames() {
		return nil, source.frameNum, fmt.Errorf("End of file")
	}

	img, err := source.MovieExtractor.ExtractFrame(source.frameNum)
	source.frameNum++

	return img, source.frameNum - 1, err
}

func (source *MovieExtractorFrameSource) FrameNum() uint64 {
	return source.frameNum
}
