package filter

import (
	"context"
	"image"
	"testing"

	"github.com/WIZARDISHUNGRY/hls-await/internal/corpus"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/sirupsen/logrus"
)

//go:generate sh -c "go test ./... -run '^$' -benchmem -bench . | tee benchresult.txt"
//go:generate sh -c "git show :./benchresult.txt | go run golang.org/x/perf/cmd/benchstat -delta-test none -geomean /dev/stdin benchresult.txt | tee benchdiff.txt"

func TestMinDistFromCorpus(t *testing.T) {
	testPatterns, err := corpus.LoadEmbedded("testpatterns")
	if err != nil {
		t.Fatalf("Load testpatterns: %v", err)
	}
	interesting, err := corpus.LoadFS("interesting")
	if err != nil {
		t.Fatalf("Load interesting: %v", err)
	}

	f := DefaultMinDistFromCorpus(testPatterns)
	ctx := testCtx()
	for name, img := range interesting.ImagesMap() {
		ok, err := f(ctx, img)
		if err != nil {
			t.Fatalf("filter: %v", err)
		}
		if !ok {
			t.Fatalf("filter failed for %s", name)
		}
	}
}

func TestMinDistFromCorpus_rejects_self(t *testing.T) {
	testPatterns, err