package main

import (
	"github.com/amarburg/go-movieset"
	"github.com/bamiaux/rez"
	"image"
	"image/png"
	"log"
	"os"
)

type ImageMaker struct {
	mm    movieset.MovieExtractor
	scale float32
	ot    OutTree
}

type Images struct {
	ImgPath, ThumbPath string
	FrameNum           uint64
}

func NewImageMaker(mm movieset.MovieExtractor, ot OutTree) *ImageMaker {
	return &ImageMaker{
		mm:    mm,
		scale: 0.1,
		ot:    ot,
	}
}

func (im *ImageMaker) MakeImages(frameNum uint64) Images {
	imgFilename := im.ot.ImageFileName(frameNum)
	thumbFilename := im.ot.ThumbnailFileName(frameNum)

	images := Images{
		ImgPath:   im.ot.RelaImageFileName(frameNum),
		ThumbPath: im.ot.RelaThumbnailFileName(frameNum),
		FrameNum:  frameNum,
	}

	_, errImg := os.Stat(imgFilename)
	_, thumbImg := os.Stat(thumbFilename)
	if errImg == nil && thumbImg == nil {
		log.Printf("Skipping %s", imgFilename)
		return images
	}

	img, _ := im.mm.ExtractFrame(frameNum)
	imgFile, _ := os.Create(imgFilename)
	png.Encode(imgFile, img)
	imgFile.Close()

	thumb := image.NewRGBA(image.Rect(0, 0,
		int(float32(img.Bounds().Dx())*im.scale),
		int(float32(img.Bounds().Dy())*im.scale)))
	rez.Convert(thumb, img, rez.NewBicubicFilter())
	thumbFile, _ := os.Create(thumbFilename)
	png.Encode(thumbFile, thumb)
	thumbFile.Close()

	return images
}
