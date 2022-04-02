package segment

import (
	"bytes"
	"encoding/gob"
	"image"
	"testing"
)

func TestGobEnc(t *testing.T) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc 