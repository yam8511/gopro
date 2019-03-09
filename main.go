package main

import (
	"context"
	"log"
	"time"
)

func main() {
	conf := ZoolConf{
		MaxCap: 1,
		Factory: func() (interface{}, error) {
			now := time.Now()
			log.Println("Factory --> ", now)
			return now, nil
		},
		Close: func(t interface{}) {
			log.Println("Closed --> ", t)
		},
		Return: func(t interface{}) interface{} {
			log.Println("Return --> ", t)
			return t
		},
		IdelTime: time.Second * 5,
	}

	zoo, err := NewZool(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Len ---> ", zoo.Len())

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		var key interface{} = "t"
		go func() {
			defer cancel()
			t, err := zoo.Get(ctx)
			if err != nil && err != ErrContextHasDone {
				log.Println("Get Error --->", err)
				return
			}
			log.Println("Now2 ---> ", t)
			ctx = context.WithValue(ctx, key, t)
		}()

		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				log.Println("Ctx Done --->", ctx.Err())
			case context.Canceled:
				t := ctx.Value(key)
				log.Println("Ctx Done --->", t)
				time.Sleep(time.Second)
				zoo.Put(t)
			}
		}
	}()

	go func() {
		t, err := zoo.Get(nil)
		if err != nil {
			log.Println("Get Error --->", err)
			return
		}
		log.Println("Now ---> ", t)
		zoo.Put(t)
	}()

	log.Println("A Len ---> ", zoo.Len())
	time.Sleep(time.Second * 6)
	log.Println("B Len ---> ", zoo.Len())

	zoo.Release()
	log.Println("C Len ---> ", zoo.Len())
}
