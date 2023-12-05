package batchpool_test

import (
	"testing"

	bp "github.com/tomr-ninja/batchpool"
)

func TestBatcher(t *testing.T) {
	n := 100
	batchSize := 10
	expectedBatches := n / batchSize

	batcher := bp.NewBatcher[testBatchElement](batchSize, expectedBatches) // all batches fit in the buffer
	k := 1
	if err := batcher.Range(func(v *testBatchElement) (next bool, err error) {
		v.value = k
		k++

		return k <= n, nil
	}); err != nil {
		t.Error(err)
	}

	batches := 0
	kMax := 0
	for b := range batcher.Batches() {
		batches++
		for _, v := range b.Data() {
			if v.value > kMax {
				kMax = v.value
			}
		}
		b.Close()
	}
	if batches != expectedBatches {
		t.Errorf("expected %d batches, got %d", expectedBatches, batches)
	}
	if kMax != n {
		t.Errorf("expected max value %d, got %d", n, kMax)
	}
}
