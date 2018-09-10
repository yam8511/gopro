package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	zz "gopro/io"
	"gopro/testdata"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// TCPAddress tcp位址
const TCPAddress = "ip:port"

// PomeloAddress Pomelo位址
const PomeloAddress = "ip:port"

func main() {
	nanoTCP()
	// peerTCP()
	// wsHands()
}

// 使用 nano 的 client tcp 進行連線
func nanoTCP() {
	addr := PomeloAddress
	c := zz.NewConnector()
	chReady := make(chan int)
	c.OnConnected(func() {
		fmt.Println("連線囉")
		chReady <- 1
	})

	fmt.Println("開始連線")
	if err := c.Start(addr); err != nil {
		fmt.Println("連線錯誤 --->", err)
		return
	}

	<-chReady

	fmt.Println("＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝")
	time.Sleep(5 * time.Second)

	// fmt.Println("PING")
	if c.IsClosed() {
		fmt.Println("連線已經關閉")
		return
	}

	data := &testdata.Enter{
		A: "Zoular",
		P: "qwe123",
		G: 11001,
	}

	c.Request("connector.entryHandler.enter", data, func(data interface{}) {
		fmt.Println("＊＊＊＊＊＊＊＊＊＊＊", data)
		return
	})

	time.Sleep(30 * time.Second)
	fmt.Println("＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝")
	defer c.Close()
}

// 單純用TCP連線
func peerTCP() {
	c, err := net.Dial("tcp", TCPAddress)
	if err != nil {
		fmt.Println("連線錯誤 --->", err)
		return
	}

	// 	s := `{
	// 	"sys": {
	// 		"version": "1.1.1",
	// 		"type": "tcp"
	// 	},
	// 	"user": {	}
	// }`
	// 	data, _ := json.Marshal(s)

	// 	fmt.Println(string(data))

	data := []byte(" Zuolar")

	buf := new(bytes.Buffer)
	buf.WriteByte(0x01)
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, len(data))

	data = append(buf.Bytes(), data...)
	_, err = c.Write(data)
	if err != nil {
		fmt.Println("寫入錯誤 --->", err)
		return
	}
	fmt.Println("寫入成功")

	_, err = c.Read(data)
	if err != nil {
		fmt.Println("讀取錯誤 --->", err)
		return
	}
	fmt.Println("讀取成功", string(data))

	// _, err = c.Write([]byte{0x01, 0, 0, 1, 1})
	// if err != nil {
	// 	fmt.Println("寫入錯誤 --->", err)
	// 	return
	// }
	// fmt.Println("寫入成功 --->", err)

	// _, err = c.Read(data)
	// if err != nil {
	// 	fmt.Println("讀取錯誤 --->", err)
	// 	return
	// }
	// fmt.Println("Hello World", string(data))
}

// 使用 ws 連線
func wsHands() {
	c, _, err := websocket.DefaultDialer.Dial("ws://"+PomeloAddress, nil)
	if err != nil {
		fmt.Println("WS 連線錯誤 --->", err)
		return
	}
	handshake := []byte{1, 0, 0, 0}
	req := map[string]interface{}{
		"sys": map[string]string{
			"version": "1.1.1",
			"type":    "js-websocket",
		},
		"user": map[string]string{},
	}
	data, _ := json.Marshal(req)
	handshake = append(handshake, data...)
	fmt.Println("握手參數 --->", string(handshake))

	err = c.WriteMessage(websocket.BinaryMessage, handshake)
	if err != nil {
		fmt.Println("WS 寫入錯誤 --->", err)
		return
	}

	go func() {
		mt, p, err := c.ReadMessage()
		if err != nil {
			fmt.Println("WS 讀取錯誤 --->", err)
			return
		}
		fmt.Println("讀取資料 型態 --->", mt)
		fmt.Println("讀取資料 --->", p)
		fmt.Println("讀取資料 --->", string(p))
	}()

}
