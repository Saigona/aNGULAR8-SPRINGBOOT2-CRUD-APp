package imagescore

import (
	"context"
	"image"
	"sync"
	"testing"
)

//go:generate sh -c "go test ./... -run '^$' -benchmem -bench . | tee benchresult.txt"
//go:generate sh -c "git show :./benchresult.txt | go run golang.org/x/perf/cmd/benchstat -delta-test none -geomean /dev/stdin benchresult.txt | tee benchdiff.txt"

func BenchmarkBulkScores(b *testing.B) {
	const (
		xDim = 720
		yDim = 576
	)
	rect := image.Rectangle{Min: image.Point{}, Max: image.Point{X: xDim, Y: yDim}}

	for _, tC := range standardTestCases {
		b.Run(tC.desc, func(b *testing.B) {
			ctx := context.Background()
			bs := NewBulkScore(ctx, tC.scoreF)
			b.ResetTimer()
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					img := image.NewRGBA(rect)
					_, err := bs.ScoreImage(ctx, img)
					if err != nil {
						b.Fatalf("ScoreImage: %v", err)
					}
				}
			})

		})
	}

}

func FuzzBulk(f *testing.F) {
	f.Fuzz(func(t *testing