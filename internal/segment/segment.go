package segment

import (
	"encoding/gob"
	"image"
)

type Handler interface {
	HandleSegment(request *Reques