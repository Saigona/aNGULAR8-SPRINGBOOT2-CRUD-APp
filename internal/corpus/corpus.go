package corpus

import (
	"embed"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/exp/maps"
)

//go:embed images/uninteresting images/testpatterns
var content embed.FS

type Corpus struct {
	name   string
	images map[string]image.Image
}

func (c *Corpus) Images() []image.Image {
	return maps.Values(c.images)
}
func (c *Corpus) ImagesMap() map[string]image.Image {
	return c.images
}
func (c *Corpus) Name() string {
	return c.name
}

var (
	_, b, _, _ = runtime.Caller(0)
	corpusRoot = filepath.Dir(b)
)

func LoadFS(path string) (*Corpus, error) {
	c := &Corpus{
		name:   path,
		images: make(map[string]image.Image),
	}

	err := filepath.Walk(
		