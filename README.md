# go-movieset

[![GoDoc](https://godoc.org/github.com/amarburg/go-movieset?status.svg)](https://godoc.org/github.com/amarburg/go-movieset)

Define two interfaces:  `MovieExtractor` which supports random access of
frames from a sequences of images (a movie or similar); and `FrameSource`
which defines a one-way iterator for a sequence of images.

It includes a set of abstractions which then sit on top off / build on these
interfaces.
