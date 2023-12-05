package batchpool_test

import (
	"math/rand"
	"testing"

	"github.com/tomr-ninja/batchpool"
)

type dataPoint struct {
	v1 int64
	v2 int64
	v3 int64
	v4 int64
}

func BenchmarkBatcher(b *testing.B) {
	batchSize := 1000

	b.Run("no batching", func(b *testing.B) {
		c := int64(0)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ps := make([]*dataPoint, 0, batchSize) // GC will have to clean this up
			for j := 0; j < batchSize; j++ {
				p := &dataPoint{} // allocated every time
				randomDataPoint(p)
				ps = append(ps, p)
			}
			c += ps[len(ps)-1].v1
		}
	})

	b.Run("batching", func(b *testing.B) {
		batcher := batchpool.NewBatcher[dataPoint](1000, 8)

		done := make(chan struct{})
		go func() {
			c := int64(0)
			for batch := range batcher.Batches() {
				ps := batch.Data()
				c += ps[len(ps)-1].v1
				batch.Close()
			}
			close(done)
		}()

		b.ReportAllocs()
		b.ResetTimer()

		i := 0
		k := 0
		_ = batcher.Range(func(v *dataPoint) (next bool, err error) {
			randomDataPoint(v)
			if k%batchSize == 0 {
				i++
			}
			k++

			return i < b.N, nil
		})
		<-done
	})
}

func randomDataPoint(p *dataPoint) {
	p.v1 = rand.Int63()
	p.v2 = rand.Int63()
	p.v3 = rand.Int63()
	p.v4 = rand.Int63()
}
