package main

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
func NewBatcher(size int, waitTime int64, numWorkers int, funct Function) (b *BatchConfig, err error) {

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
		WaitTime:  time.Duration(waitTime) * time.Second,
		Func:      funct,
		items:     make([]interface{}, 0),
		batchChan: make(chan []interface{}, numWorkers),
	}
	batch.worker(numWorkers)

	go batch.autoDump()

	go batch.timeout()

	return batch, nil
}
//This function helps to insert item to array
func (b *BatchConfig) Insert(item interface{}) bool{
	if len(b.items) < b.size {
		mutex.Lock()
		b.items = append(b.items, item)
		mutex.Unlock()
		return true
	}
	return false
}

//This function helps to auto dump when the items reach size.
func (b *BatchConfig) autoDump(){
	if len(b.items) == b.size {
		mutex.Lock()
		b.dump()
		mutex.Unlock()
	}
}

func (b *BatchConfig) dump() {
	data := make([]interface{}, len(b.items))

	for id, item := range b.items {
		data[id] = item
	}
	b.items = b.items[:0]
	b.batchChan <- data
}

//This function heps to with timeout dump
func (b *BatchConfig) timeout() {
	ticker := time.NewTicker(b.WaitTime)
	for {
		select {
		case <-ticker.C:
			mutex.Lock()
			if len(b.items) != 0  {
				b.dump()
			}
			mutex.Unlock()
		}
	}
}
