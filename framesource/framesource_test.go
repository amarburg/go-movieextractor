package framesource

import (
	"github.com/amarburg/go-lazyfs-testfiles/frameset"
	"testing"
)

func TestFrameSourceGoodJson(t *testing.T) {

	source, err := MakeFrameSourceFromPath(frameset_testfiles.GoodMultiMovJson)

	if err != nil {
		t.Errorf("Unable to make frame source from good.json: %s", err)
	}

	goodMultiMovJsonTest( t, source )


}
