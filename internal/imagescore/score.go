package imagescore

import (
	"context"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"io"
	"sync/atomic"

	"github.com/WIZARDISHUNGRY/h