package frameset

import (
	"github.com/amarburg/go-lazyfs-testfiles/frameset"
	"testing"
)

func TestLoadGoodJson(t *testing.T) {

	set, err := LoadFrameSet(frameset_testfiles.GoodFrameSetJson)

	if err != nil {
		t.Errorf("Unable to load good.json: %s", err)
	}

	const NumChunks = frameset_testfiles.GoodFrameSetJsonChunks
	FrameLengths := []int{3, 0, 2, 0, 0}
	ChunkMins := []uint64{1, 20, 30, 50, 60}
	ChunkStarts := []uint64{0, 20, 0, 50, 60}
	ChunkEnds := []uint64{0, 30, 0, 60, 65}

	if len(set.Chunks) != NumChunks {
		t.Errorf("Unexpected number of chunks %d != %d", len(set.Chunks), NumChunks)
	}

	for i, chunk := range set.Chunks {

		if len(chunk.Frames) != FrameLengths[i] {
			t.Errorf("Unexpected frame length in chunk %d; %d != %d", i, len(chunk.Frames), FrameLengths[i])
		}

		if chunk.Min() != ChunkMins[i] {
			t.Errorf("Unexpected chunk min in chunk %d; %d != %d", i, chunk.Min(), ChunkMins[i])
		}

		if chunk.Start != ChunkStarts[i] {
			t.Errorf("Unexpected chunk start in chunk %d; %d != %d", i, chunk.Min(), ChunkStarts[i])
		}

		if chunk.End != ChunkEnds[i] {
			t.Errorf("Unexpected chunk end in chunk %d; %d != %d", i, chunk.Min(), ChunkEnds[i])
		}

	}

}
