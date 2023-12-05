## Usage

```go
type DataPoint struct {
    Time   time.Time
    Values map[string]float64
}

func main() {
    const (
        ssssssreadWorkers    = 10
        uploadWorkers  = 10
        batchSize      = 100
        batchesBufSize = 10
    )

    batcher := batchpool.NewBatcher[DataPoint](batchSize, batchesBufSize)

    var wg sync.WaitGroup
    wg.Add(readWorkers + uploadWorkers)

    for i := 0; i < readWorkers; i++ {
        go func() {
            defer wg.Done()
            
            for f := range filesToRead() {
                defer f.Close()

                s := bufio.NewScanner(f)

                if err := batcher.Range(func(v *DataPoint) (next bool, err error) {
                    if !s.Scan() {
                        return false, nil
                    }

                    return true, parseLine(v, s.Bytes())
                }); err != nil {
                    panic(err)
                }
            }
        }()
    }

    for i := 0; i < uploadWorkers; i++ {
        go func() {
            defer wg.Done()

            for b := range batcher.Batches() {
                if err := uploadDataPoints(b.Data()); err != nil {
                    panic(err)
                }

                b.Close()
            }
        }()
    }
}

func filesToRead() <-chan *os.File {
	panic("not implemented")
}

func parseLine(_ *DataPoint, _ []byte) error {
	panic("not implemented")
}

func uploadDataPoints(_ []*DataPoint) error {
	panic("not implemented")
}
```