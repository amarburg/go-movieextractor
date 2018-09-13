package movieset

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

type NotAFrameSetError struct{}

func (f NotAFrameSetError) Error() string {
	return "Not a frameset"
}

// LoadFrameSet reads a FrameSet...
func LoadFrameSet(path string) (*FrameSet, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	set := new(FrameSet)
	err = decoder.Decode(set)

	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		return nil, err
	}

	if len(set.Source) == 0 {
		return nil, NotAFrameSetError{}
	}

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

func (fs *FrameSet) ChunkNames() []string {
	var names []string
	for _, chunk := range fs.Chunks {
		names = append(names, chunk.Name)
	}
	return names
}

func (fs *FrameSet) MovFromChunk(name string) (VirtualMov, error) {
	for _, chunk := range fs.Chunks {

		if chunk.Name == name {

			ex, err := fs.MovieExtractor()

			if err != nil {
				return VirtualMov{}, err
			}

			return CreateVirtualMov(ex, chunk.Start, chunk.End)
		}
	}

	return VirtualMov{}, fmt.Errorf("Couldn't find chunk \"%s\"", name)
}

func (set FrameSet) MovieExtractor() (MovieExtractor, error) {
	// Create the source

	source := os.ExpandEnv(set.Source)

	// If path is relative, evaluate it relative to the input file...
	if !filepath.IsAbs(source) {
		if len(set.filepath) == 0 {
			return nil, fmt.Errorf("Source path in frameset is relative, but don't know filepath for frameset")
		}

		source = filepath.Clean(filepath.Join(filepath.Dir(set.filepath), source))
	}

	return OpenMovieExtractor(source)

}
