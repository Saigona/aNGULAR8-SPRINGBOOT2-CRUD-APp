package segment

import (
	"bytes"
	"encoding/gob"
	"image"
	"testing"
)

func TestGobEnc(t *testing.T) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read fr