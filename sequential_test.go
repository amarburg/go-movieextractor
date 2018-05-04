package movieset

import (
	"github.com/amarburg/go-lazyfs-testfiles/frameset"
	"github.com/amarburg/go-lazyfs-testfiles/multimov"
	"testing"
	"io"
)

func goodSequentialCheck(t *testing.T, source Sequential) {
	frames := 0
	done := false

	for done == false {
		_, _, err := source.Next()

		switch err {
		case io.EOF:
			done = true
		case nil:
		default:
			t.Fatalf("Error reading frame: %s", err)
		}

		frames++
		t.Logf("Got frame %d", frames)

		// TODO.  Check that frames are valid.
	}

	if frames != frameset_testfiles.GoodFrameSetJsonFrames {
		t.Errorf("Didn't get as many frames as I expected %d, rather than %d",
							frameset_testfiles.GoodFrameSetJsonFrames, frames)
	}
}


func TestSequentialGoodJson(t *testing.T) {

	source, err := OpenSequential(frameset_testfiles.GoodFrameSetJson)

	if err != nil {
		t.Errorf("Unable to make frame source from good.json: %s", err)
	}

	goodSequentialCheck(t, source)

}

func TestFrameSourceMultimov(t *testing.T) {

	source, err := OpenSequential(multimov_testfiles.FourMovMultiMovJson)

	if err != nil {
		t.Errorf("Unable to make frame source from %s: %s", multimov_testfiles.FourMovMultiMovJson, err)
	}

	_, frame, err := source.Next()

	if err != nil {
		t.Errorf("Unable to retrieve image from %s: %s", multimov_testfiles.FourMovMultiMovJson, err)
	}

	if frame != 1 {
		t.Errorf("Didn't get frame 1 from %s, got %d", multimov_testfiles.FourMovMultiMovJson, frame)
	}
}
