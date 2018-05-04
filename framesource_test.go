package movieset

import (
	"github.com/amarburg/go-lazyfs-testfiles/frameset"
	"github.com/amarburg/go-lazyfs-testfiles/multimov"
	"testing"
)

func TestFrameSourceGoodJson(t *testing.T) {

	source, err := MakeFrameSourceFromPath(frameset_testfiles.GoodMultiMovJson)

	if err != nil {
		t.Errorf("Unable to make frame source from good.json: %s", err)
	}

	goodFrameSourceTest(t, source)

}

func TestFrameSourceMultimov(t *testing.T) {

	source, err := MakeFrameSourceFromPath(multimov_testfiles.FourMovMultiMovJson)

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
