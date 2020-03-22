package batcher

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func (b *BatchConfig) worker (workerCount int){
	for i := 1; i < workerCount; i++ {
		go func(workerID int, flushJobs <-chan []interface{}) {
			for j := range flushJobs {
				fmt.Println(j)
				b.Func(workerID,j)
			}
		}(i, b.batchChan)

	}
}




