package multimov

// A SequenceElement represents one movie within a sequence
type SequenceElement struct {
	FrameOffset uint64
	Hash        MovHash
}

// Sequence is a convenience type representing a slice of SequenceElements
type Sequence []SequenceElement

// func (seq *Sequence) ToJson(dest io.Writer) error {
// 	encoder := json.NewEncoder(dest)
//
// 	return encoder.Encode(seq)
// }
//
// func SequenceFromJson(src io.Reader) (*Sequence, error) {
// 	decoder := json.NewDecoder(src)
//
// 	seq := NewSequence()
//
// 	err := decoder.Decode(seq)
//
// 	return seq, err
// }
//
// func NewSequence() *Sequence {
// 	return &Sequence{
// 		movies: make([]MovRecord, 0),
// 	}
// }
