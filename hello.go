package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gopkg.in/robfig/cron.v2"
)

// CronJob 排程背景
type CronJob struct {
	// 背景名稱
	Name string
	// 執行週期
	Spec string
	// 執行工作
	Cmd func()
	// 是否可以重複
	IsOverlapping bool
	// EntryID
	EntryID cron.EntryID
	// running
	running bool
	// 讀寫鎖
	mux *sync.RWMutex
	// 等待通道
	wg *sync.WaitGroup
}

// Run 執行背景
func (c *CronJob) Run() {
	log.Printf("開始執行 ------> %s\n", c.Name)

	c.mux.RLock()
	overlapping := c.IsOverlapping
	running := c.running
	c.mux.RUnlock()

	// 如果可以重複，直接執行
	if overlapping {
		c.wg.Add(1)
		c.Cmd()
		c.wg.Done()
		return
	}

	// 如果不可重複，而且已經執行則跳過
	if running {
		log.Printf("還在執行中 ------> %s\n", c.Name)
		return
	}

	// 執行背景
	c.wg.Add(1)
	c.mux.Lock()
	c.running = true
	c.mux.Unlock()
	c.Cmd()
	c.mux.Lock()
	c.running = false
	c.mux.Unlock()
	c.wg.Done()
}

// Init 初始化
func (c *CronJob) Init() {
	c.mux = new(sync.RWMutex)
	c.wg = new(sync.WaitGroup)
}

// Wait 等待結束
func (c *CronJob) Wait() {
	c.wg.Wait()
}

func main() {
	fmt.Println("Hello World")
	c := cron.New()

	jobs := []*CronJob{
		{Name: "echo_a", Spec: "*/2 * * * * *", Cmd: echoa, IsOverlapping: true},  // 可重複執行的任務
		{Name: "echo_z", Spec: "*/2 * * * * *", Cmd: echoz, IsOverlapping: false}, // 不可重複執行的任務
		{Name: "sleep", Spec: "*/2 * * * * *", Cmd: keep, IsOverlapping: false},   // 不會結束的任務
	}

	for _, job := range jobs {
		job.Init()
		pid, err := c.AddJob(job.Spec, job)
		if err != nil {
			log.Fatalln(err)
		}
		job.EntryID = pid
	}

	c.Start()
	hasShutdown := false
	sig := gracefulShutdown()
	go func() {
		for {
			<-sig
			if hasShutdown {
				log.Println(`

				強制結束囉~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

				`)
				os.Exit(137)
			}
		}
	}()
	<-sig
	hasShutdown = true

	log.Println(`

	收到訊號囉~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	`)
	c.Stop()

	for _, job := range jobs {
		job.Wait()
	}

	log.Println(`

	結束囉~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	`)
}

func echoa() {
	log.Println("============= BEGIN 可重複 ==============")
	<-time.After(time.Second * 3)
	log.Println("=============  END  可重複 ==============")
}

func echoz() {
	log.Println("xxxxxxxxxxxxx BEGIN 不重複 xxxxxxxxxxxxxx")
	<-time.After(time.Second * 8)
	log.Println("xxxxxxxxxxxxx  END  不重複 xxxxxxxxxxxxxx")
}

func keep() {
	for range time.NewTicker(time.Second * 5).C {
		log.Println("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	}
}

func gracefulShutdown() (sig chan os.Signal) {
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	return
}
