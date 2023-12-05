

## Usage

```go
import bp "github.com/tomr-ninja/batchpool"

type DataPoint struct {
	Time time.Time
	Value float64
}

func main() {
	batcher := bp.NewBatcher[DataPoint](batchSize, expectedBatches)
	batcher.
	
}

```

```go

```