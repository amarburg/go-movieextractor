package movieset

import (
	"encoding/json"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyfs-testfiles"
	"github.com/amarburg/go-lazyfs-testfiles/multimov"
	"github.com/amarburg/go-lazyquicktime"
	"os"
	"testing"
)

//=== Utility functions ==

func expectMovLength(t *testing.T, mov *MultiMov, expected int) {
	if mov.NumMovies() != expected {
		t.Fatalf("Expected MultiMov to have %d movies, actually had %d", expected, mov.NumMovies())
	}
}

func expectNumFrames(t *testing.T, mov *MultiMov, expected uint64) {
	if mov.NumFrames() != expected {
		t.Fatalf("Expected MultiMov to have %d frames, actually had %d", expected, mov.NumFrames())
	}
}

func expectOffset(t *testing.T, mov *MultiMov, offset uint64, expectedOffset uint64, expectedMovie MovHash) {
	hash, offset, err := mov.Offset(offset)

	if err != nil {
		t.Fatalf("Error looking up offset of %d: %s", offset, err)
	}

	if offset != expectedOffset {
		t.Fatalf("Expected offset of %d, got %d", expectedOffset, offset)
	}

	if hash != expectedMovie {
		t.Fatalf("Expected movie %x, got %x", expectedMovie, hash)
	}
}

func expectOffsetError(t *testing.T, mov *MultiMov, offset uint64) {
	_, _, err := mov.Offset(offset)

	if err == nil {
		t.Fatalf("Expected an error when looking up offset %d, didn't get one", offset)
	}
}

func expectOffsetPath(t *testing.T, mov *MultiMov, offset uint64, expectedPath string) {
	hash, _, _ := mov.Offset(offset)
	movPath := mov.MovPath(hash)
	if movPath != expectedPath {
		t.Fatalf("Expected movie path \"%s\", not \"%s\"", expectedPath, movPath)
	}
}

//== Tests ==

func TestZeroConstructor(t *testing.T) {
	var mm MultiMov

	if mm.NumMovies() != 0 {
		t.Error("null constructed MultiMov isn't length 0")
	}

	expectMovLength(t, &mm, 0)
}

func TestJsonUnmarshalEmptyMultiMov(t *testing.T) {
	mm, err := LoadMultiMov(multimov_testfiles.EmptyMultiMovJson)

	if err != nil {
		t.Fatalf("Error while loading empty multimov test file %s: %s", multimov_testfiles.EmptyMultiMovJson, err)
	}

	expectMovLength(t, mm, 0)
}

func TestJsonUnmarshalZeroLengthMultiMov(t *testing.T) {
	mm, err := LoadMultiMov(multimov_testfiles.ZeroLengthMultiMovJson)

	if err != nil {
		t.Fatalf("Error while loading zero length multimov test file %s: %s", multimov_testfiles.ZeroLengthMultiMovJson, err)
	}

	expectMovLength(t, mm, 0)

	expectOffsetError(t, mm, 0)
	expectOffsetError(t, mm, 1)
}

func TestLoadOneMultiMov(t *testing.T) {
	mm, err := LoadMultiMov(multimov_testfiles.SingleMovMultiMovJson)

	if err != nil {
		t.Fatalf("Error while loading single entry multimov test file %s: %s", multimov_testfiles.SingleMovMultiMovJson, err)
	}

	expectMovLength(t, mm, 1)
	expectNumFrames(t, mm, uint64(lazyfs_testfiles.TestMovNumFrames))

	// Try frame lookup
	expectOffset(t, mm, 1, 0, mm.Sequence[0].Hash)
	expectOffset(t, mm, uint64(lazyfs_testfiles.TestMovNumFrames), 0, mm.Sequence[0].Hash)

	expectOffsetError(t, mm, 0)
	expectOffsetError(t, mm, uint64(lazyfs_testfiles.TestMovNumFrames)+1)

}

func TestLoadFourMultiMov(t *testing.T) {
	mm, err := LoadMultiMov(multimov_testfiles.FourMovMultiMovJson)

	if err != nil {
		t.Fatalf("Error while loading single entry multimov test file %s: %s", multimov_testfiles.SingleMovMultiMovJson, err)
	}

	expectMovLength(t, mm, 4)
	expectNumFrames(t, mm, 4*uint64(lazyfs_testfiles.TestMovNumFrames))

	expectOffset(t, mm, 1, 0, mm.Sequence[0].Hash)
	expectOffset(t, mm, uint64(lazyfs_testfiles.TestMovNumFrames), 0, mm.Sequence[0].Hash)

	expectOffsetError(t, mm, 0)
	expectOffsetError(t, mm, (4*uint64(lazyfs_testfiles.TestMovNumFrames))+1)

	expectOffsetPath(t, mm, 1, lazyfs_testfiles.TestMovPath)

}

//== Test appending a movie to an empty

func TestAppendMovie(t *testing.T) {
	mm := NewMultiMov()

	fs, err := lazyfs.OpenLocalFile(lazyfs_testfiles.TestMovPath)
	if err != nil {
		t.Fatalf("Error opening local lazyfs: %s", err)
	}

	expectMovLength(t, &mm, 0)

	lqt, err := lazyquicktime.LoadMovMetadata(fs)
	if err != nil {
		t.Fatalf("Error loading quicktime movie %s: %s", lazyfs_testfiles.TestMovPath, err)
	}

	mm.AppendMovie(MovRecordFromLqt(lqt))

	expectMovLength(t, &mm, 1)
	expectNumFrames(t, &mm, uint64(lazyfs_testfiles.TestMovNumFrames))

	mm.AppendMovie(MovRecordFromLqt(lqt))

	expectMovLength(t, &mm, 2)
	expectNumFrames(t, &mm, 2*uint64(lazyfs_testfiles.TestMovNumFrames))

	file, err := os.Create("temp.json")

	if err != nil {
		t.Fatalf("Failed to create temporary file \"temp.json\"")
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(mm)

	if err != nil {
		t.Fatalf("Encoder failed")
	}

}

//
//
// func TestSequence(t *testing.T) {
// 	seq := NewSequence()
//
// 	if seq.movies == nil {
// 		t.Error("seq.records is unset")
// 	}
//
// }
//
// func TestJson(t *testing.T) {
// 	seq := NewSequence()
//
// 	buf := bytes.NewBufferString("")
//
// 	err := seq.ToJson(buf)
//
// 	if err != nil {
// 		t.Error(fmt.Sprintf("Error serializing to JSON: %s", err.Error()))
// 	}
//
// 	unser, err := SequenceFromJson(buf)
//
// 	if err != nil {
// 		t.Error(fmt.Sprintf("Error unserializing from JSON: %s", err.Error()))
// 	}
//
// 	if len(unser.movies) != len(seq.movies) {
// 		t.Error(fmt.Sprintf("Sequences different length after unserialization %d != %d", len(unser.movies), len(seq.movies)))
// 	}
//
// }
