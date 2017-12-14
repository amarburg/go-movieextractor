package main

import (
	"github.com/amarburg/go-lazyquicktime"
	"github.com/bamiaux/rez"
	"image"
	"image/png"
	"log"
	"os"
)

type ImageMaker struct {
	mm    lazyquicktime.MovieExtractor
	scale float32
	ot    OutTree
}

type Images struct {
	ImgPath, ThumbPath         string
	RelaImgPath, RelaThumbPath string
	FrameNum                   uint64
}

func NewImageMaker(mm lazyquicktime.MovieExtractor, ot OutTree) *ImageMaker {
	return &ImageMaker{
		mm:    mm,
		scale: 0.25,
		ot:    ot,
	}
}

func (im *ImageMaker) MakeImages(frameNum uint64) Images {
	imgFilename := im.ot.ImageFileName(frameNum)
	thumbFilename := im.ot.ThumbnailFileName(frameNum)

	images := Images{
		RelaImgPath:   im.ot.RelaImageFileName(frameNum),
		RelaThumbPath: im.ot.RelaThumbnailFileName(frameNum),
		ImgPath:       im.ot.ImageFileName(frameNum),
		ThumbPath:     im.ot.ThumbnailFileName(frameNum),
		FrameNum:      frameNum,
	}

	var img image.Image

	_, errErr := os.Stat(imgFilename)
	if errErr != nil {
		log.Printf("Making %s", imgFilename)

		img, _ = im.mm.ExtractFrame(frameNum)
		imgFile, _ := os.Create(imgFilename)
		png.Encode(imgFile, img)
		imgFile.Close()
	} else {
		log.Printf("Skipping %s", imgFilename)
	}

	_, thumbErr := os.Stat(thumbFilename)
	if thumbErr != nil {
		log.Printf("Making thumbnail %s", thumbFilename)

		if img == nil {
			imgFile, err := os.Create(imgFilename)
			if err != nil {
				log.Printf("Error opening %s: %s", imgFilename, err)
			}
			img, err = png.Decode(imgFile)
			if err != nil {
				log.Printf("Error decoding png %s: %s", imgFilename, err)
			 }
			imgFile.Close()
		}

		if img != nil {
		thumb := image.NewRGBA(image.Rect(0, 0,
			int(float32(img.Bounds().Dx())*im.scale),
			int(float32(img.Bounds().Dy())*im.scale)))
		rez.Convert(thumb, img, rez.NewBicubicFilter())
		thumbFile, _ := os.Create(thumbFilename)
		png.Encode(thumbFile, thumb)
		thumbFile.Close()
}
	}

	return images
}
