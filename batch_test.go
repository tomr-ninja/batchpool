package batchpool_test

import (
	"testing"

	bp "github.com/tomr-ninja/batchpool"
)

type testBatchElement struct {
	value int
}

func TestBatch(t *testing.T) {
	isNotZero := func(v *testBatchElement) bool {
		return v.value != 0
	}
	isZero := func(v *testBatchElement) bool {
		return v.value == 0
	}

	doNothingAfterReset := func(_ *bp.Batch[testBatchElement]) {}

	t.Run("range once", func(t *testing.T) {
		batchSize := 10

		b := bp.NewBatch[testBatchElement](batchSize, doNothingAfterReset)

		k := 1
		cb := func(v *testBatchElement) (next bool, err error) {
			v.value = k
			k++

			return k <= batchSize, nil
		}

		if err := b.Range(cb); err != nil {
			t.Error(err)
		}

		allSet := false
		_ = b.Range(func(v *testBatchElement) (next bool, err error) {
			allSet = isNotZero(v)

			return true, nil
		})
		if !allSet {
			t.Error("expected all values to be set")
		}

		allReset := false
		b.Close()
		_ = b.Range(func(v *testBatchElement) (next bool, err error) {
			allReset = isZero(v)

			return true, nil
		})
		if !allReset {
			t.Error("expected all values to be reset")
		}
	})
}
