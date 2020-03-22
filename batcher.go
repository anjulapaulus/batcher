package batcher

import (
	"errors"
	"sync"
	"time"
)

type Function func(workerID int, data []interface{})

type BatchConfig struct {
	size      int
	WaitTime  time.Duration
	Func      Function
	items     []interface{}
	batchChan chan []interface{}
}

var mutex = &sync.Mutex{}

//Initialises a new instance
func NewBatcher(size int, waitTime time.Duration, numWorkers int, funct Function) (b *BatchConfig, err error) {

	switch {
	case size <= 0:
		return &BatchConfig{}, errors.New("invalid size")
	case waitTime <= 0:
		return &BatchConfig{}, errors.New("invalid wait time")
	case numWorkers <= 0:
		return &BatchConfig{}, errors.New("invalid number of workers")
	}

	batch := &BatchConfig{
		size:      size,
		WaitTime:  waitTime,
		Func:      funct,
		items:     make([]interface{}, 0), //initialize empty slice
		batchChan: make(chan []interface{}, numWorkers),
	}
	batch.worker(numWorkers)

	go batch.autoDump()
	go batch.timeout()

	return batch, nil
}

//This function helps to insert item to array
func (b *BatchConfig) Insert(item interface{}) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if len(b.items) < b.size {
		b.items = append(b.items, item)
		return true
	}
	//if len(b.items) == b.size {
	//	b.dump()
	//}

	return false
}


func (b *BatchConfig) dump() {
	copiedItems := make([]interface{},len(b.items))
	if len(b.items) != 0 {
		mutex.Lock()
		copy(copiedItems,b.items)
		b.items = b.items[:0]
		b.batchChan <- copiedItems
		mutex.Unlock()
	}
}

func (b *BatchConfig) autoDump() {
	if len(b.items) == b.size {
		b.dump()
	}
}

//This function heps to with timeout dump
func (b *BatchConfig) timeout() {
	ticker := time.NewTicker(b.WaitTime)
	for {
		select {
		case <- ticker.C:
				b.dump()
		}
	}
}
