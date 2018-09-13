package movieset

import (
	"fmt"
	"image"
	"time"
)

type VirtualMov struct {
	Offset, Length uint64
	Mov            MovieExtractor
}

func CreateVirtualMov(mov MovieExtractor, offset uint64, length uint64) (VirtualMov, error) {
	// l := uint64(length)
	// if length < 0 {
	//   l = mov.NumFrames() - uint64(offset)
	// }

	return VirtualMov{
		Mov:    mov,
		Offset: offset,
		Length: length,
	}, nil

}

func (vm VirtualMov) NumFrames() uint64 {
	return vm.Length
}

func (vm VirtualMov) Duration() time.Duration {
	// TODO
	return 0
}

func (vm VirtualMov) ExtractFrame(frame uint64) (image.Image, error) {
	vframe := frame + vm.Offset
	if vframe >= vm.Length {
		return nil, fmt.Errorf("Frame %d is beyond virtual length %d", vframe, vm.Length)
	}
	return vm.Mov.ExtractFrame(vframe)
}

//ExtractFrame(frame uint64) (image.Image, error)
//ExtractFramePerf(frame uint64) (image.Image, LQTPerformance, error)
