package io

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"

	"gopro/codec"
	"gopro/message"
	"gopro/packet"

	"github.com/golang/protobuf/proto"
)

var (
	hsd []byte // handshake data
	had []byte // handshake ack data
	hbt []byte // heartbeat data
)

func init() {
	req := map[string]interface{}{
		"sys": map[string]string{
			"version": "1.1.1",
			"type":    "js-websocket",
		},
		// "user": map[string]string{},
	}
	data, _ := json.Marshal(req)

	var err error
	hsd, err = codec.Encode(packet.Handshake, data)
	if err != nil {
		panic(err)
	}

	had, err = codec.Encode(packet.HandshakeAck, nil)
	if err != nil {
		panic(err)
	}

	hbt, err = codec.Encode(packet.Heartbeat, nil)
	if err != nil {
		panic(err)
	}
}

type (

	// Callback represents the callback type which will be called
	// when the correspond events is occurred.
	Callback func(data interface{})

	// Connector is a tiny Nano client
	Connector struct {
		conn   net.Conn       // low-level connection
		codec  *codec.Decoder // decoder
		die    chan struct{}  // connector close channel
		chSend chan []byte    // send queue
		mid    uint           // message id

		// events handler
		muEvents sync.RWMutex
		events   map[string]Callback

		// response handler
		muResponses sync.RWMutex
		responses   map[uint]Callback

		connectedCallback func() // connected callback
		close             bool
	}
)

// NewConnector create a new Connector
func NewConnector() *Connector {
	return &Connector{
		die:       make(chan struct{}),
		codec:     codec.NewDecoder(),
		chSend:    make(chan []byte, 64),
		mid:       1,
		events:    map[string]Callback{},
		responses: map[uint]Callback{},
	}
}

// Start connect to the server and send/recv between the c/s
func (c *Connector) Start(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

	go c.write()

	// send handshake packet
	c.send(hsd)

	// read and process network message
	go c.read()

	return nil
}

// OnConnected set the callback which will be called when the client connected to the server
func (c *Connector) OnConnected(callback func()) {
	c.connectedCallback = callback
}

// Request send a request to server and register a callbck for the response
func (c *Connector) Request(route string, v proto.Message, callback Callback) error {
	data, err := serialize(v)
	if err != nil {
		return err
	}

	log.Println("請求 ---> ", string(data))

	msg := &message.Message{
		Type:  message.Request,
		Route: route,
		ID:    c.mid,
		Data:  data,
	}

	c.setResponseHandler(c.mid, callback)
	if err := c.sendMessage(msg); err != nil {
		c.setResponseHandler(c.mid, nil)
		return err
	}

	return nil
}

// Notify send a notification to server
func (c *Connector) Notify(route string, v proto.Message) error {
	data, err := serialize(v)
	if err != nil {
		return err
	}

	msg := &message.Message{
		Type:  message.Notify,
		Route: route,
		Data:  data,
	}
	return c.sendMessage(msg)
}

// On add the callback for the event
func (c *Connector) On(event string, callback Callback) {
	c.muEvents.Lock()
	defer c.muEvents.Unlock()

	c.events[event] = callback
}

// Close close the connection, and shutdown the benchmark
func (c *Connector) Close() {
	if c.close {
		return
	}
	log.Println("連線關閉")
	c.conn.Close()
	close(c.die)
	c.close = true
}

// IsClosed check the connection is closed
func (c *Connector) IsClosed() bool {
	return c.close
}

func (c *Connector) eventHandler(event string) (Callback, bool) {
	c.muEvents.RLock()
	defer c.muEvents.RUnlock()

	cb, ok := c.events[event]
	return cb, ok
}

func (c *Connector) responseHandler(mid uint) (Callback, bool) {
	c.muResponses.RLock()
	defer c.muResponses.RUnlock()

	cb, ok := c.responses[mid]
	return cb, ok
}

func (c *Connector) setResponseHandler(mid uint, cb Callback) {
	c.muResponses.Lock()
	defer c.muResponses.Unlock()

	if cb == nil {
		delete(c.responses, mid)
	} else {
		c.responses[mid] = cb
	}
}

func (c *Connector) sendMessage(msg *message.Message) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}

	// log.Printf("資料 ---> %+v,\n整體 ---> %+v\n加密 ---> %+v\n", msg.Data, msg, data)

	payload, err := codec.Encode(packet.Data, data)
	if err != nil {
		return err
	}

	c.mid++
	c.send(payload)

	return nil
}

func (c *Connector) write() {
	for {
		select {
		case data := <-c.chSend:
			log.Println("準備傳送訊息 --->", data)
			if _, err := c.conn.Write(data); err != nil {
				log.Println("傳送訊息失敗", err.Error())
				// c.Close()
			}

		case <-c.die:
			return
		}
	}
}

func (c *Connector) send(data []byte) {
	c.chSend <- data
}

func (c *Connector) read() {
	buf := make([]byte, 2048)

	for {
		time.Sleep(time.Second)
		if c.IsClosed() {
			return
		}
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Println("讀取資料失敗", err.Error())
			// c.Close()
			// return
			continue
		}

		packets, err := c.codec.Decode(buf[:n])
		if err != nil {
			log.Println("解碼資料失敗", err.Error())
			// c.Close()
			// return
			continue
		}

		for i := range packets {
			p := packets[i]
			log.Println("讀取到資料包 --->", p)
			c.processPacket(p)
		}
	}
}

func (c *Connector) processPacket(p *packet.Packet) {
	switch p.Type {
	case packet.Handshake:
		var handShakeResponse struct {
			Code int `json:"code"`
			Sys  struct {
				Heartbeat int `json:"heartbeat"`
			} `json:"sys"`
			// User struct{} `json:"user"`
		}

		err := json.Unmarshal(p.Data, &handShakeResponse)
		if err != nil {
			log.Fatal("握手回傳進行解碼發生錯誤")
			c.Close()
			return
		}

		if handShakeResponse.Code == 200 {
			go func() {
				ticker := time.NewTicker(time.Second * time.Duration(handShakeResponse.Sys.Heartbeat))
				for range ticker.C {
					if c.IsClosed() {
						log.Println("沒心跳了，不要發")
						return
					}
					log.Println("發送心跳包給主子")
					c.send(hbt)
				}
			}()
			c.send(had)
			c.connectedCallback()
		} else {
			log.Fatal("握手回傳不是200狀態", string(p.Data))
			c.Close()
		}
	case packet.Data:
		msg, err := message.Decode(p.Data)
		if err != nil {
			log.Println(err.Error())
			return
		}
		c.processMessage(msg)

	case packet.Kick:
		log.Fatal("Server 主動斷開連線通知 --->", p)
		c.Close()
	}
}

func (c *Connector) processMessage(msg *message.Message) {
	switch msg.Type {
	case message.Push:
		cb, ok := c.eventHandler(msg.Route)
		if !ok {
			log.Println("event handler not found", msg.Route)
			return
		}

		cb(msg.Data)

	case message.Response:
		cb, ok := c.responseHandler(msg.ID)
		if !ok {
			log.Println("response handler not found", msg.ID)
			return
		}

		cb(msg.Data)
		c.setResponseHandler(msg.ID, nil)
	}
}

func serialize(v proto.Message) ([]byte, error) {
	data, err := proto.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}
