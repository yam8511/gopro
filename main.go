package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	c := flag.Bool("c", false, "客戶端")
	n := flag.Int("n", 1, "連線數量")
	d := flag.Int("d", 1, "呼叫數量")
	flag.Parse()

	if !*c {
		runServer()
	}

	count := *n
	done := make(chan error, count)
	for i := 0; i < count; i++ {
		go func() {
			err := runClient(*d)
			done <- err
		}()
		if i%99 == 0 {
			time.Sleep(time.Millisecond * 500)
		}
	}

	msg := map[string]int{}
	ok := 0
	for i := 0; i < count; i++ {
		err := <-done
		if err != nil {
			c, _ := msg[err.Error()]
			c++
			msg[err.Error()] = c
		} else {
			ok++
		}
	}

	for k, v := range msg {
		log.Println(k, v)
	}
	log.Println("OK,", ok)
}

func runClient(n int) error {
	url := "http://127.0.0.1:8000/api"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("New Request : " + err.Error())
	}

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "26d7fa44-e531-4d7f-8ec0-1f27810efd2a")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Do CLient : " + err.Error())
	}

	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Read Body : " + err.Error())
	}
	return nil
}

func NewListenerZ(l net.Listener, max int) ListenerZ {
	z := ListenerZ{
		l:       l,
		maxConn: max,
		in:      make(chan int),
		out:     make(chan int),
		get:     make(chan int),
	}
	go z.kernal()
	return z
}

type ListenerZ struct {
	l           net.Listener
	maxConn     int
	currentConn int
	in          chan int
	out         chan int
	get         chan int
}

func (z ListenerZ) Accept() (net.Conn, error) {
	// log.Println("***** Acecept ")
	z.getPool()
	// log.Println("***** Acecept Done")

	conn, err := z.l.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (z ListenerZ) Close() error {
	return z.l.Close()
}

func (z ListenerZ) Addr() net.Addr {
	return z.l.Addr()
}

func (z ListenerZ) kernal() {
	for {
		select {
		case <-z.out:
			// 有人離開請求
			z.currentConn--
			// log.Println("有人離開請求,", z.currentConn)
		case z.in <- 1:
			// 有人來取請求
			z.currentConn++
			// log.Println("有人來取請求,", z.currentConn)
		case z.get <- z.currentConn:
			// 有人來確認目前數量
		}
		time.Sleep(time.Millisecond)
	}
}

func (z ListenerZ) getPool() {
	for {
		select {
		case currentCount := <-z.get:
			if currentCount < z.maxConn {
				<-z.in
				return
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (z ListenerZ) putPool() {
	z.out <- 0
}

func runServer() {
	r := gin.New()
	r.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "ok, "+c.Request.RemoteAddr+", "+time.Now().Format(time.RFC3339Nano))
	})
	// r.Run(":8000")

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	// zl := NewListenerZ(l, 1023)
	server := &http.Server{
		Handler: r,
		// ConnState: func(conn net.Conn, state http.ConnState) {
		// 	log.Printf("Conn: %+v, State: %v\n", conn.RemoteAddr(), state)
		// 	if state == http.StateClosed {
		// 		zl.putPool()
		// 	}
		// },
	}
	// server := &http.Server{
	// 	Handler: r,
	// }

	server.SetKeepAlivesEnabled(false)
	// err = serve(server, l)
	// err = server.Serve(zl)
	err = server.Serve(l)
	log.Fatal(err)
}
