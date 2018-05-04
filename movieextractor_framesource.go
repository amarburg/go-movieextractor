package movieset

import (
	"fmt"
	"image"
)

// Thin wrapper around a MovieExtractor which implements FrameSource
type MovieExtractorFrameSource struct {
	ext      MovieExtractor
	frameNum uint64
}

// func (ext MovieExtractor) FrameSource() (*MovieExtractorFrameSource, error) {
// 	return &MovieExtractorFrameSource{
// 		MovieExtractor: ext,
// 		frameNum:       1,
// 	}, nil
// }

func FrameSourceFromMovieExtractor(ext MovieExtractor) (*MovieExtractorFrameSource, error) {

	return &MovieExtractorFrameSource{
		ext:      ext,
		frameNum: 1,
	}, nil
}

func (source *MovieExtractorFrameSource) Next() (image.Image, uint64, error) {

	if source.frameNum > source.ext.NumFrames() {
		return nil, source.frameNum, fmt.Errorf("End of file")
	}

	img, err := source.ext.ExtractFrame(source.frameNum)
	source.frameNum++

	return img, source.frameNum - 1, err
}

func (source *MovieExtractorFrameSource) FrameNum() uint64 {
	return source.frameNum
}
