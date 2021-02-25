package bot

import (
	"fmt"
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

func TestParseDate(t *testing.T) {
	tw := &twitter.Tweet{
		CreatedAt: `Wed Feb 23 23:25:53 +0000 2022`,
	}
	t