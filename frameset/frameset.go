package frameset

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type FrameSet struct {
	filepath  string
	Source    string
	Chunks    SliceOfChunks `json:",omitempty"`
	ImageName string
	NumFrames uint64 `json:",omitempty"`
}

type SliceOfChunks []Chunk
type UInt64Slice []uint64

type Chunk struct {
	Name        string
	Description string      `json:",omitempty"`
	Start       uint64      `json:",omitempty"`
	End         uint64      `json:",omitempty"`
	Frames      UInt64Slice `json:",omitempty"`

	//Chunks SliceOfChunks `json:",omitempty"`
}

func (c Chunk) Min() uint64 {
	if c.HasFrames() {
		// Assume sorted
		return c.Frames[0]
	}

	return c.Start
}

func (c Chunk) HasFrames() bool {
	return len(c.Frames) > 0
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

	// First validate.   Can either have Start or Frames but not both
	for i, chunk := range set.Chunks {
		if chunk.HasFrames() && chunk.Start != 0 {
			return nil, fmt.Errorf("Chunk \"%s\" has both frames and a start", chunk.Name)
		}

		// Sort frames if presento
		sort.Sort(set.Chunks[i].Frames)
	}

	// Ensure the chunks are sorted
	sort.Sort(set.Chunks)

	// Link chunks
	for i, chunk := range set.Chunks {
		if chunk.Start != 0 && chunk.End == 0 {

			if i == len(set.Chunks)-1 {
				set.Chunks[i].End = set.NumFrames
			} else {
				set.Chunks[i].End = set.Chunks[i+1].Min()
			}
		}
	}

	set.filepath = path

	return set, err
}

// Makes SliceOfChunks sortable
func (p SliceOfChunks) Len() int           { return len(p) }
func (p SliceOfChunks) Less(i, j int) bool { return p[i].Min() < p[j].Min() }
func (p SliceOfChunks) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p UInt64Slice) Len() int           { return len(p) }
func (p UInt64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p UInt64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
