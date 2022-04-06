package segment

import (
	"encoding/gob"
	"image"
)

type Handler interface {
	HandleSegment(request *Request, resp *Response) error // yes, an interface pointer as first arg, we'll 