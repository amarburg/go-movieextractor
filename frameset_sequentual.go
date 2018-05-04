package movieset

import (
	"fmt"
	"image"
	"io"
)

type FrameSetSequential struct {
	*FrameSet
	Movie         MovieExtractor
	chunkIdx      int
	frameIdx      int
	segmentOffset uint64
	totalFrames   uint64
}

func MakeFrameSetSequential(set *FrameSet) (*FrameSetSequential, error) {

	mm, err := set.MovieExtractor()

	if err != nil {
		return &FrameSetSequential{}, err
	}

	return &FrameSetSequential{
		FrameSet: set,
		Movie:    mm,
	}, nil
}

func (source *FrameSetSequential) Valid() error {
	if source.chunkIdx >= len(source.FrameSet.Chunks) {
		return io.EOF
	}

	chunk := source.FrameSet.Chunks[source.chunkIdx]

	if chunk.HasFrames() {
		if source.frameIdx >= len(chunk.Frames) {
			return fmt.Errorf("Frame offset is off end of frame array (error) in chunk %d; %d >= %d", source.chunkIdx, source.frameIdx, len(chunk.Frames))
		}
	} else if (chunk.Start + source.segmentOffset) >= chunk.End {
		return fmt.Errorf("Segment offset is off end of segment (error) in chunk %d; %d >= %d", source.chunkIdx, (chunk.Start + source.segmentOffset), chunk.End)
	}

	return nil
}

func (source *FrameSetSequential) Advance() {
	source.frameIdx++
	source.segmentOffset++
	source.totalFrames++

	chunk := source.FrameSet.Chunks[source.chunkIdx]

	if chunk.HasFrames() {

		if source.frameIdx >= len(chunk.Frames) {
			source.frameIdx = 0
			source.segmentOffset = 0
			source.chunkIdx++
		}
	} else if (chunk.Start + source.segmentOffset) >= chunk.End {
		source.frameIdx = 0
		source.segmentOffset = 0
		source.chunkIdx++
	}

}

func (source *FrameSetSequential) Next() (image.Image, uint64, error) {
	if err := source.Valid(); err != nil {
		return nil, 0, err
	}

	chunk := source.FrameSet.Chunks[source.chunkIdx]

	var frame uint64
	if chunk.HasFrames() {
		frame = chunk.Frames[source.frameIdx]
	} else {
		frame = chunk.Start + source.segmentOffset
	}

	defer source.Advance()

	img, err := source.Movie.ExtractFrame(frame)
	return img, frame, err
}

// func (source *FrameSetSequential) FrameNum() uint64 {
// 	return source.totalFrames
// }
