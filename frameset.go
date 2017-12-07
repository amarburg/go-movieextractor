package frameset

import (
	"encoding/json"
	//"log"
	"os"
	"sort"
)

type FrameSet struct {
	Source    string
	Frames    []uint64
	Chunks    ChunkSlice `json:",omitempty"`
	ImageName string
	NumFrames uint64 `json:",omitempty"`
}

type ChunkSlice []Chunk

type Chunk struct {
	Name  string
	Start uint64
	End		uint64 `json:",omitempty"`
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

	// Ensure the chunks are sorted
	sort.Sort(set.Chunks)

	// Link sprites

	for i, _ := range set.Chunks {
		if set.Chunks[i].End == 0 {

		if i == len(set.Chunks)-1 {
			set.Chunks[i].End = set.NumFrames
		} else {
			set.Chunks[i].End = set.Chunks[i+1].Start
		}
	}
}

	return set, err
}



// Makes ChunkSlice sortable
func (p ChunkSlice) Len() int           { return len(p) }
func (p ChunkSlice) Less(i, j int) bool { return p[i].Start < p[j].Start }
func (p ChunkSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
