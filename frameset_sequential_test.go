package movieset

import (
	"github.com/amarburg/go-lazyfs-testfiles/frameset"
	"testing"
)



func TestFrameSetSequentialGoodJson(t *testing.T) {

	set, err := LoadFrameSet(frameset_testfiles.GoodFrameSetJson)

	if err != nil {
		t.Errorf("Unable to load good.json: %s", err)
	}

	source, err := MakeFrameSetSequential(set)

	if err != nil {
		t.Errorf("Unable to make frame source from good.json: %s", err)
	}

	goodSequentialCheck(t, source)

}
