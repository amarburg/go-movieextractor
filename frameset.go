package multimov

import (
	"encoding/json"
	"os"
	"sort"
)

type FrameSet struct {
	Source    string
	Frames    []uint64
	Chunks    map[string]Chunk `json:",omitempty"`
	ImageName string
	NumFrames uint64 `json:",omitempty"`
}

type Chunk struct {
	Range FrameRange
}

type FrameRange struct {
	Start uint64
	End   uint64 `json:",omitempty"`
}

// LoadMultiMov reads a MultiMov from the a path to a given JSON file.
// Returns a pointer to a new MultiMov if successful, or nil and
// an error if unsuccessful
func LoadFrameSet(path string) (*FrameSet, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	set := new(FrameSet)
	err = decoder.Decode(set)

	return set, err
}

type NamedChunk struct {
	Chunk
	Key string
}

type NamedChunkSlice []NamedChunk

func (p NamedChunkSlice) Len() int           { return len(p) }
func (p NamedChunkSlice) Less(i, j int) bool { return p[i].Range.Start < p[j].Range.Start }
func (p NamedChunkSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (set FrameSet) OrderedChunks() NamedChunkSlice {
	out := make(NamedChunkSlice, 0, len(set.Chunks))

	for key, chunk := range set.Chunks {
		out = append(out, NamedChunk{Chunk: chunk, Key: key})
	}

	sort.Sort(out)

	// Fill out end ranges
	for i := 0; i < len(out); i++ {
		if out[i].Range.End == 0 {

			if i == len(out)-1 {
				out[i].Range.End = set.NumFrames
			} else {
				out[i].Range.End = out[i+1].Range.Start
			}
		}
	}

	return out
}
