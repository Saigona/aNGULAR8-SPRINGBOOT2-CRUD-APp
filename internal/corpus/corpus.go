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
	ret