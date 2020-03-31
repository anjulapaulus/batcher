package batcher

import (
	"errors"
	"sync"
	"time"
)

type Function func(data []interface{}) bool

type BatchConfig struct {
	MaxCapacity       int
	WaitTime   time.Duration
	function   Function
	batchChan  chan interface{}
	mutex      sync.RWMutex
}

//Initialises a new instance
func NewBatcher(maxJobs int, waitTime time.Duration, f Function) (*BatchConfig, error) {
	switch {
	case maxJobs <= 0:
		return &BatchConfig{}, errors.New("invalid size")
	case waitTime <= 0:
		return &BatchConfig{}, errors.New("invalid wait time")
	}


	return &BatchConfig{
		MaxCapacity:  maxJobs,
		WaitTime: waitTime,
		function: f,
	}, nil
}

//This function helps to insert an item
func (b *BatchConfig) Insert(item interface{}) (bool, error) {
	if item == nil {
		return false, errors.New("item inserted is null")
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.batchChan == nil {
		b.batchChan = make(chan interface{}, b.MaxCapacity)
		go b.dumper()
	}

	b.batchChan <- item

	return true, nil
}

func (b *BatchConfig) InsertItems(items []interface{}) (bool, error) {
	if items == nil {
		return false, errors.New("items inserted is null")
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	batchLen := len(items)
	// If the length of batch is larger than maxCapacity
	if batchLen > b.MaxCapacity {
		items = items[:b.MaxCapacity]
	}

	if b.batchChan == nil {
		b.batchChan = make(chan interface{}, b.MaxCapacity)
		go b.dumper()
	}

	for _, item := range items {
		b.batchChan <- item
	}
	return true, nil
}




func (b *BatchConfig) dumper() {
	var batch []interface{}
	timer := time.NewTimer(b.WaitTime)

	for {
		select {
		case <-timer.C:
			b.function(batch)
			b.close()
			return
		case item := <-b.batchChan:
			batch = append(batch, item)
			if len(batch) >= b.MaxCapacity {
				// Callback with batch
				b.function(batch)
				// Init batch array
				batch = []interface{}{}
			}
		}
	}
}

func (b *BatchConfig) close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.batchChan != nil {
		close(b.batchChan)
		b.batchChan = nil
	}
}
