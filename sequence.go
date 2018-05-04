package movieset

// A SequenceElement represents one movie within a sequence
type SequenceElement struct {
	FrameOffset uint64
	Hash        MovHash
}

// Sequence is a convenience type representing a slice of SequenceElements
type Sequence []SequenceElement
