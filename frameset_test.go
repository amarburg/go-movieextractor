package frameset

import (
	"testing"
)

func TestLoadGoodJson(t *testing.T) {

	set, err := LoadFrameSet("test_data/good.json")

	if err != nil {
		t.Error("Unable to load good.json: %s", err)
	}

	const NumChunks = 5
	FrameLengths := []int{3, 0, 2, 0, 0}
	ChunkMins := []uint64{1, 20, 30, 50, 60}
	ChunkStarts := []uint64{0, 20, 0, 50, 60}
	ChunkEnds := []uint64{0, 30, 0, 60, 0}

	if len(set.Chunks) != NumChunks {
		t.Error("Unexpected number of chunks %d != %d", len(set.Chunks), NumChunks)
	}

	for i, chunk := range set.Chunks {

		if len(chunk.Frames) != FrameLengths[i] {
			t.Error("Unexpected frame length %d != %d", len(chunk.Frames), FrameLengths[i])
		}

		if chunk.Min() != ChunkMins[i] {
			t.Error("Unexpected chunk min %d != %d", chunk.Min(), ChunkMins[i])
		}

		if chunk.Start != ChunkStarts[i] {
			t.Error("Unexpected chunk start %d != %d", chunk.Min(), ChunkStarts[i])
		}

		if chunk.End != ChunkEnds[i] {
			t.Error("Unexpected chunk end %d != %d", chunk.Min(), ChunkEnds[i])
		}

	}

}
