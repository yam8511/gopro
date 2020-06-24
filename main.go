package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

const N = 10

var familyCtx = [N]context.Context{}
var familyCancel = [N]context.CancelFunc{}

func main() {
	wg := sync.WaitGroup{}
	ch := make(chan int)
	for i := 0; i < N; i++ {
		familyCtx[i], familyCancel[i] = context.WithCancel(context.Background())
		wg.Add(1)
		go func(i int) {
			rand.Seed(time.Now().UnixNano())
			d := time.Duration(rand.Intn(10)+2) * time.Second
			log.Println("No.", i, ", 睡覺:", d)
			wg.Done()
			time.Sleep(d)
			familyCancel[i]()
			ch <- i
		}(i)
	}
	wg.Wait()

	for i := 0; i < N-1; i++ {
		log.Println("=============")
		n := <-ch
		familyCtx[n] = nil
		log.Printf("第%d輪, No. %d 起床\n", i, n)
	}
	log.Println("=============")

	for i := 0; i < N; i++ {
		if familyCtx[i] != nil {
			log.Println("睡到最後: No.", i)
		}
	}
}
