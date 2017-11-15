package FrameSet


type FrameSet struct {
	Source    string
	Frames    []uint64
	ImageName string
  Chunks    map[string][]struct {
              Frames []uint64
            }
}
