# Batcher
Library Implementation for Batch Processing in GO.

## Installation
````
go get github.com/anjulapaulus/batcher
````

## Implementation

Single Insert

````
import (
	"fmt"
	"github.com/anjulapaulus/batcher"
	"log"
	"time"
)



func main() {
	b, err := batcher.NewBatcher(1000, 5*time.Millisecond, doBatch)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for i := 1; i <= 1000; i++ {
			b.Insert(i)
		}
	}()

	fmt.Println("Stopping")
	time.Sleep(time.Second * 60)
}

func doBatch(datas []interface{}) bool {

	for _, data := range datas {
		if parsedValue, ok := data.(int); ok {
			log.Println(fmt.Sprintf("[data]: %d", parsedValue))
		}
	}
	return true
}
````

Bulk Insert

````
func main() {
	arr:=[]interface{}{1,2,3}
	b, err := batcher.NewBatcher(1000, 5*time.Millisecond, doBatch)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
			b.InsertItems(arr)
	}()

	fmt.Println("Stopping")
	time.Sleep(time.Second * 5)
}

func doBatch(datas []interface{}) bool {

	for _, data := range datas {
		if parsedValue, ok := data.(int); ok {
			log.Println(fmt.Sprintf("[data]: %d", parsedValue))
		}
	}
	return true
}

````