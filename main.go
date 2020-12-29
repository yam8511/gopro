package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"gopkg.in/fsnotify.v1"
)

const NODES = "./data/nodes"
const PUBLIC = "./data/public"

var ips []string
var ipmx *sync.RWMutex
var getIPs func() []string
var setIPs func(newIPs []string)
var ipCH = make(chan []string, 1)

func init() {
	// ====== STEP 0 ======
	// >>>>>> 讀取所有node
	nodeBytes, err := ioutil.ReadFile(NODES)
	if err != nil {
		log.Fatal(err)
	}

	ips = parseNodes(nodeBytes)
	ipmx = &sync.RWMutex{}
	getIPs = func() []string {
		ipmx.RLock()
		defer ipmx.RUnlock()
		return ips
	}
	setIPs = func(newIPs []string) {
		ipmx.Lock()
		ips = newIPs
		ipmx.Unlock()
	}

	// >>>>>> 監控node檔案更動
	go watchNodes()
}

func main() {
	/**
	0. 先知道所有node的IP [讀取nodes, 開啟瀏覽器]
	1. 加入 Registry IP 到每台node的`/etc/hosts`
	2. 指定node成為 Master
	3. 指定node成為 Server and Agent
	4.
	*/

	// >>>>>> 開啟瀏覽器
	var port string
	for {
		port = "8000"
		fmt.Print("請輸入要開啟瀏覽器的Port [8000 - 35565] 不輸入預設為8000 \n> ")
		fmt.Scanln(&port)
		port = strings.TrimSpace(port)
		p, err := strconv.Atoi(port)
		if err != nil || p < 8000 || p > 35565 {
			log.Println("請輸入數字 [8000 - 35565]")
		} else {
			break
		}
	}

	gin.SetMode(gin.ReleaseMode)
	var r *gin.Engine
	if isDebug() {
		r = gin.Default()
	} else {
		r = gin.New()
		r.Use(gin.Recovery())
	}
	r.NoRoute(static.ServeRoot("/", PUBLIC))
	deploy := r.Group("/deploy")
	{
		deploy.GET("/guide", guide)
		deploy.POST("/registry", func(c *gin.Context) {})
		deploy.POST("/master", func(c *gin.Context) {})
		deploy.POST("/server", func(c *gin.Context) {})
		deploy.POST("/monitoring", func(c *gin.Context) {})
		deploy.POST("/logging", func(c *gin.Context) {})
		deploy.POST("/dashboard", func(c *gin.Context) {})
	}

	if !isDebug() {
		browser.OpenURL("http://127.0.0.1:" + port)
	}

	log.Println("伺服器監聽: " + port)
	err := r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func isDebug() bool {
	return os.Getenv("DEBUG") == "1"
}

func inNodes(node string, nodes []string) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}

	return false
}

func parseNodes(nodeBytes []byte) []string {
	ips := []string{}
	nodes := strings.Split(strings.TrimSpace(string(nodeBytes)), "\n")
	for _, n := range nodes {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		ips = append(ips, n)
	}
	return ips
}
func watchNodes() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	watcher.Add(NODES)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if isDebug() {
				log.Println("event =>", event)
			}
			switch event.Op {
			case fsnotify.Write, fsnotify.Chmod:
				nodeBytes, err := ioutil.ReadFile(event.Name)
				if err != nil {
					log.Fatal(err)
				}
				newIPs := parseNodes(nodeBytes)
				setIPs(newIPs)
				select {
				case ipCH <- newIPs:
				case <-time.After(time.Second):
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
