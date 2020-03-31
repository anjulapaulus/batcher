package batcher

import (
	"testing"
	"time"
)



func DummyBatchFn1(data []interface{}) bool{
	return true
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
		 _,err :=NewBatcher(tt.batch.size, tt.batch.waitTime,tt.batch.funct)
		 if err!= nil{
		 	if tt.response != false{
		 		t.Error("[ERROR]: NewBatcher function")
			}
		 }

	}

}


func TestBatchConfig_Insert(t *testing.T) {
	batch,err := NewBatcher(60, 10*time.Millisecond, DummyBatchFn1)

	if err != nil{
		t.Error("[ERROR]: NewBatcher function - Insert")
	}

	for i:=1; i<=1000; i++ {
		_, err :=batch.Insert(i)
		if err != nil{
			t.Error("[ERROR]: Insert function")
		}

	}
}


func TestBatchConfig_InsertItems(t *testing.T) {
	batch, err := NewBatcher(10, 10*time.Millisecond, DummyBatchFn1)
	arr:=[]interface{}{1,2,3}
	if err != nil{
		t.Error("[ERROR]: NewBatcher function - Insert")
	}
	_, err =batch.InsertItems(arr)
	if err != nil{
		t.Error("[ERROR]: Insert function")
	}
}

