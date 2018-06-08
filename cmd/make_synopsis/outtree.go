package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type OutTree interface {
	ThumbnailFileName(idx uint64) string
	ImageFileName(idx uint64) string
	RelaThumbnailFileName(idx uint64) string
	RelaImageFileName(idx uint64) string
}

type DirOutTree struct {
	basedir string
}

func NewDirOutTree(basedir string) DirOutTree {
	_ = os.MkdirAll(basedir, os.ModePerm)

	return DirOutTree{
		basedir: basedir,
	}
}

func (ot DirOutTree) join(p string) string {
	return filepath.Join(ot.basedir, p)
}

func (ot DirOutTree) thumbDir() string {
	path := filepath.Join(ot.basedir, "thumbs")

	// eager-create directory
	_ = os.MkdirAll(path, os.ModePerm)

	return path
}

func (ot DirOutTree) imageDir() string {
	path := filepath.Join(ot.basedir, "images")

	// eager-create directory
	_ = os.MkdirAll(path, os.ModePerm)

	return path
}

func (ot DirOutTree) ThumbnailFileName(idx uint64) string {
	return filepath.Join(ot.thumbDir(), fmt.Sprintf("thumb_%08d.png", idx))
}

func (ot DirOutTree) ImageFileName(idx uint64) string {
	return filepath.Join(ot.imageDir(), fmt.Sprintf("frame_%08d.png", idx))
}

func (ot DirOutTree) RelaThumbnailFileName(idx uint64) string {
	return filepath.Join("thumbs", fmt.Sprintf("thumb_%08d.png", idx))
}

func (ot DirOutTree) RelaImageFileName(idx uint64) string {
	return filepath.Join("images", fmt.Sprintf("frame_%08d.png", idx))
}
