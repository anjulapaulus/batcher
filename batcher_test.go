package batcher

import (
	"fmt"
	"log"
	"testing"
)

func DummyBatchFn1(workerID int, data []interface{}) {
	fmt.Println(workerID, data)
	for _,v := range data{
		log.Println(v)
	}
}

type batch struct {
	size       int
	waitTime   int64
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
			waitTime:   1,
			numWorkers: 1,
			funct:      DummyBatchFn1,
		}, false,
		},
		{
			batch{
				size:       1,
				waitTime:   0,
				numWorkers: 1,
				funct:      DummyBatchFn1,
			}, false,
		},
		{
			batch{
				size:       1,
				waitTime:   1,
				numWorkers: 0,
				funct:      DummyBatchFn1,
			}, false,
		},
		{
			batch{
				size:       1,
				waitTime:   15,
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
	batch, err := NewBatcher(2,60,1, DummyBatchFn1)
	if err !=nil{
		t.Error("Failed: Insert Function Test : New Batcher")
	}
	for i:=1; i<=2; i++ {
		insert := batch.Insert(i)
		if insert != true {
			t.Error("Failed: Insert Function Test")
		}
	}

}
