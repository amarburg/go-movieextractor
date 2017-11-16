package FrameSet

type FrameSet struct {
	Source    string
	Frames    []uint64
	ImageName string
	Chunks map[string][]uint64
	// Chunks    map[string]struct {
	// 	Frames []uint64
	// }
}
