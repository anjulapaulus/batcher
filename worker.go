package main

import (
	"sync"
)

var wg sync.WaitGroup

func (b *BatchConfig) worker (workerCount int){
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerId int, workers <- chan []interface{}){
			for work := range workers{
				b.Func(i,work)
			}
			wg.Done()
		}(i, b.batchChan)

	}
}



