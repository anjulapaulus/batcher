package batcher

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func DummyBatchFn1(workerID int, data []interface{}) {
	for _,v  := range  data{
		log.Println(fmt.Sprintf("[WokerID]: %d [data]: %d",workerID, v))
	}
}

type batch struct {
	size       int
	waitTime   time.Duration
	numWorkers int
	funct      Function
}

func TestNewBatcher(t *testing.T) {
	newBatcherTest := []struct {
		batch    batch
		response bool
	}{
		{batch{
			size:       0,
			waitTime:   1*time.Second,
			numWorkers: 1,
			funct:      DummyBatchFn1,
		}, false,
		},
		{
			batch{
				size:       1,
				waitTime:   0*time.Second,
				numWorkers: 1,
				funct:      DummyBatchFn1,
			}, false,
		},
		{
			batch{
				size:       1,
				waitTime:   1*time.Second,
				numWorkers: 0,
				funct:      DummyBatchFn1,
			}, false,
		},
		{
			batch{
				size:       1,
				waitTime:   15*time.Second,
				numWorkers: 1,
				funct:      DummyBatchFn1,
			}, true,
		},
	}

	for _, tt := range newBatcherTest {
		_, err := NewBatcher(tt.batch.size, tt.batch.waitTime, tt.batch.numWorkers, tt.batch.funct)

		if err != nil {
			response := false
			if response != tt.response {
				t.Error("Failed: NewBatcher Test")
			}
		} else {
			response := true
			if response != tt.response {
				t.Error("Failed: NewBatcher Test")
			}
		}
	}

}

func TestBatchConfig_Insert(t *testing.T) {
	batch, err := NewBatcher(10,3,2, DummyBatchFn1)
	if err !=nil{
		t.Error("Failed: Insert Function Test : New Batcher")
	}
	for i:=1; i<=10; i++ {
		check := batch.Insert(i)

		if check != true {
			t.Error("Failed: Insert Function Test")
		}
	}

}
