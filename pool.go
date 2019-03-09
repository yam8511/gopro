package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

// 連接持相關錯誤
var (
	ErrConfigNotProvide = errors.New("MaxCap, Factory and Close Func must be provided")
	ErrContextHasDone   = errors.New("Context Has Done")
	ErrPoolHasClosed    = errors.New("Pool Has Closed")
)

// Zool 連線池
type Zool struct {
	conf         ZoolConf       // 池設定
	mx           sync.RWMutex   // 讀寫鎖
	idleConnPool []*ZoolConn    // 閒置中的池子
	usedConnNum  int            // 使用中的連線數量
	channel      chan *ZoolConn // 通道
	closed       bool           // 關閉狀態
}

// ZoolConf 連線池設定
type ZoolConf struct {
	MaxCap   int                           // 最大空間
	Factory  func() (interface{}, error)   // 連線工廠
	Close    func(interface{})             // 關閉連線的Func
	Return   func(interface{}) interface{} // 連線歸還後的Func
	IdelTime time.Duration                 // 閒置時間
}

// ZoolConn 連線
type ZoolConn struct {
	Conn      interface{}
	CreatedAt time.Time
}

// NewZool 產生新連線池
func NewZool(conf ZoolConf) (*Zool, error) {
	if conf.MaxCap == 0 || conf.Factory == nil || conf.Close == nil {
		return nil, ErrConfigNotProvide
	}

	zoo := &Zool{
		conf:         conf,
		idleConnPool: make([]*ZoolConn, 0),
		channel:      make(chan *ZoolConn),
	}

	// 開始中樞運作
	go func() {
		ticker := time.NewTicker(zoo.conf.IdelTime)
		for now := range ticker.C {
			// 移除超過閒置時間的
			zoo.mx.Lock()

			if zoo.closed {
				zoo.mx.Unlock()
				return
			}

			rm := []int{}
			for i := range zoo.idleConnPool {
				zooConn := zoo.idleConnPool[i]
				if zooConn.CreatedAt.Add(zoo.conf.IdelTime).Before(now) {
					go zoo.conf.Close(zooConn.Conn)
					rm = append(rm, i)
				}
			}

			for i := range rm {
				index := rm[i]
				index -= i
				zoo.idleConnPool = append(
					zoo.idleConnPool[:index],
					zoo.idleConnPool[index+1:]...,
				)
			}

			zoo.mx.Unlock()
		}
	}()

	return zoo, nil
}

// Len 目前池子數量
func (zz *Zool) Len() int {
	zz.mx.RLock()
	n := len(zz.idleConnPool) + zz.usedConnNum
	zz.mx.RUnlock()
	return n
}

// Release 釋放連線池
func (zz *Zool) Release() {
	zz.mx.Lock()
	for i := range zz.idleConnPool {
		zooConn := zz.idleConnPool[i]
		if zooConn != nil {
			zz.conf.Close(zooConn.Conn)
		}
	}
	zz.idleConnPool = []*ZoolConn{}
	zz.mx.Unlock()
}

// Put 歸還連線
func (zz *Zool) Put(conn interface{}) {
	zzLen := zz.Len()
	if zzLen == 0 {
		return
	}

	// 檢查歸還連線後是否需要執行Func
	returnFunc := zz.conf.Return
	if returnFunc != nil {
		conn = returnFunc(conn)
	}

	zz.mx.Lock()

	zooConn := &ZoolConn{Conn: conn}
	select {
	case zz.channel <- zooConn:
		// 丟到通道看看，如果有人正在等待取連線，就直接取走
	default:
		// 丟到閒置的連線池
		zooConn.CreatedAt = time.Now()
		zz.usedConnNum--
		zz.idleConnPool = append(zz.idleConnPool, zooConn)
	}

	zz.mx.Unlock()

	return
}

// Get 取連線
func (zz *Zool) Get(ctx context.Context) (interface{}, error) {
	//  檢查目前連線數量
	// ，如果連線足夠，直接建立連線
	// 如果連線超過最大數量，等待其他連線歸還

	zz.mx.Lock()

	if zz.closed {
		zz.mx.Unlock()
		return nil, ErrPoolHasClosed
	}

	idleConnNum := len(zz.idleConnPool)
	if idleConnNum > 0 {
		zooConn := zz.idleConnPool[0]
		zz.idleConnPool = zz.idleConnPool[1:]
		zz.usedConnNum++

		if zooConn != nil {
			zz.mx.Unlock()
			return zooConn.Conn, nil
		}

		conn, err := zz.conf.Factory()
		if err != nil {
			zz.usedConnNum--
		}
		zz.mx.Unlock()
		return conn, err
	}

	if zz.usedConnNum < zz.conf.MaxCap {
		zz.usedConnNum++
		conn, err := zz.conf.Factory()
		if err != nil {
			zz.usedConnNum--
		}
		zz.mx.Unlock()
		return conn, err
	}

	zz.mx.Unlock()

	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ErrContextHasDone
		case zooConn := <-zz.channel:
			if zooConn == nil {
				conn, err := zz.conf.Factory()
				if err != nil {
					zz.usedConnNum--
				}
				return conn, err
			}
			return zooConn.Conn, nil
		}
	}

	select {
	case zooConn := <-zz.channel:
		if zooConn == nil {
			conn, err := zz.conf.Factory()
			if err != nil {
				zz.usedConnNum--
			}
			return conn, err
		}
		return zooConn.Conn, nil
	}
}

// Closing 關閉中
func (zz *Zool) Closing() {
	// 先把狀態改成關閉
	zz.mx.Lock()
	zz.closed = true
	usedConnNum := zz.usedConnNum
	zz.mx.Unlock()

	// 等待連線歸還，然後關閉連線
	for i := 0; i < usedConnNum; i++ {
		zooConn := <-zz.channel
		if zooConn != nil {
			zz.conf.Close(zooConn.Conn)
		}
	}
	zz.Release()
}
