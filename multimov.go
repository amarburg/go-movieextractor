package movieset

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
)

// MultiMovVersion is the semantic version string for the current MultiMov
// JSON structure.
var MultiMovVersion = "0.1"

// A MovHash is a 32-bit hash handle to a given movie
type MovHash uint32

// MultiMov is the top-level container representing a MultiMov
type MultiMov struct {
	Version  string
	BaseDir  string `json:",omitempty"`
	Movies   map[MovHash]MovRecord
	Sequence Sequence
}

// NewMultiMov instantiates a new MultiMov
func NewMultiMov() MultiMov {
	return MultiMov{
		Version: MultiMovVersion,
		Movies:  make(map[MovHash]MovRecord),
	}
}

// LoadMultiMov reads a MultiMov from the a path to a given JSON file.
// Returns a pointer to a new MultiMov if successful, or nil and
// an error if unsuccessful
func LoadMultiMov(path string) (*MultiMov, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	mm := new(MultiMov)
	err = decoder.Decode(mm)

	if mm.BaseDir == "" {
		mm.BaseDir = filepath.Dir(path)
	}

	return mm, err
}

// Save marshals the MultiMov to a JSON file at the given path
// func (mm *MultiMov) Save(path string) error {
//
// 	// Encode MultiMov to JSON
// 	file, err := os.Create(path)
// 	defer file.Close()
//
// 	if err != nil {
// 		return nil
// 	}
//
// 	encoder := json.NewEncoder(file)
// 	encoder.SetIndent("", "  ")
// 	return encoder.Encode(mm)
// }

//== Functions for adding/removing movies from MultiMov
func (mm *MultiMov) addHash(mov MovRecord) MovHash {
	h := fnv.New32()
	io.WriteString(h, mov.ShortName)
	crc := MovHash(h.Sum32())

	mm.Movies[crc] = mov
	return crc
}

// AppendMovie adds the given MovRecord to the end of the sequence in the
// MultiMov
func (mm *MultiMov) AppendMovie(mov MovRecord) {
	crc := mm.addHash(mov)

	// Calculate offset
	var offset uint64
	if len(mm.Sequence) > 0 {
		prev := mm.Sequence[len(mm.Sequence)-1]
		offset = prev.FrameOffset + uint64(mm.Movies[prev.Hash].NumFrames)
	}

	mm.Sequence = append(mm.Sequence, SequenceElement{FrameOffset: offset, Hash: crc})
}

// NumMovies lists the number of in the MultiMov
func (mm MultiMov) NumMovies() int {
	return len(mm.Sequence)
}

// NumFrames returns the total number of frames in a MultiMov -- the sum
// of the number of frames in all movies within the MultiMov
func (mm MultiMov) NumFrames() uint64 {
	var out uint64
	for _, h := range mm.Sequence {
		out += uint64(mm.Movies[h.Hash].NumFrames)
	}

	return out
}

// Offset gives the movie handle (hash) and offset given a frame.
// The hash is absolute position of the 1st frame in the movie ...
// it's up to the user to subtract the offset from frame.
func (mm MultiMov) Offset(frame uint64) (MovHash, uint64, error) {
	for _, h := range mm.Sequence {
		mov, has := mm.Movies[h.Hash]
		if has == false {
			return 0, 0, fmt.Errorf("Error loading hash from movie table")
		}

		if frame > h.FrameOffset && frame <= h.FrameOffset+uint64(mov.NumFrames) {
			return h.Hash, h.FrameOffset, nil
		}
	}

	return 0, 0, fmt.Errorf("Couldn't identify which movie frame %d occurs in", frame)
}

// MovPath converts a movie handle to a movie file path
func (mm MultiMov) MovPath(hash MovHash) string {
	mov, has := mm.Movies[hash]

	if has == false {
		return ""
	}

	if mov.Relapath == "" {
		return filepath.Clean(filepath.Join(mm.BaseDir, mov.ShortName))
	}

	return filepath.Clean(filepath.Join(mm.BaseDir, mov.Relapath))

}
