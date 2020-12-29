package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

type Payload struct {
	Event string
	UUID  string
	Nodes []string
	Step  int
}

type Flow struct {
	Title      string
	More       bool
	NoAction   bool
	Master     bool
	Deploy     func() error
	DeployNode func(node string) error
	DeployMore func(serverNodes, agentNodes []string) error
}

func deployRegistry(node string) error {
	nodes := getIPs()
	errCh := make(chan error)
	for _, ip := range nodes {
		go func(ip string) {

			errCh <- nil
		}(ip)
	}

	txt := ""
	for i := 0; i < len(nodes); i++ {
		err := <-errCh
		if err != nil {
			txt += err.Error() + "\n"
		}
	}

	if txt != "" {
		return errors.New(txt)
	}

	return nil
}

func guide(c *gin.Context) {
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	ctx, done := context.WithCancel(context.Background())
	defer func() {
		done()
		conn.Close()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ips := <-ipCH:
				conn.WriteJSON(Payload{
					Event: "get_ip",
					Nodes: ips,
					Step:  -1,
				})
			}
		}
	}()

	var nextUUID, masterNode string
	NextUID := func() string {
		nextUUID = uuid.New().String()
		return nextUUID
	}

	flows := []Flow{
		{
			Title:      "請選擇一台 Registry Node",
			DeployNode: nil,
		},
		{
			Title:      "請選擇一台 Master Node",
			DeployNode: nil,
			Master:     true,
		},
		{
			Title:      "請選擇 Server Node (未選作為 Agent Node)",
			More:       true,
			DeployMore: nil,
		},
		{
			Title:      "請選擇一台 Node 作為日誌儲存",
			DeployNode: nil,
		},
		{
			Title:      "請選擇一台 Node 作為監控儲存",
			DeployNode: nil,
		},
		{
			Title:      "請選擇一台 Node 作為監控介面(Dashboard)",
			DeployNode: nil,
		},
		{
			Title:    "部署Database",
			NoAction: true,
			Deploy:   nil,
		},
		{
			Title:    "部署AI",
			NoAction: true,
			Deploy:   nil,
		},
		{
			Title:    "部署排程",
			NoAction: true,
			Deploy:   nil,
		},
		{
			Title:    "部署服務",
			NoAction: true,
			Deploy:   nil,
		},
	}

	for i, flow := range flows {
		step := i + 1

		log.Printf("%d. %s", step, flow.Title)

		if flow.NoAction {
			err = conn.WriteJSON(Payload{
				Event: "start",
				Step:  step,
			})
			if err != nil {
				return
			}
			flow.Deploy()
			time.Sleep(time.Second)
			err = conn.WriteJSON(Payload{
				Event: "done",
				Step:  step,
			})
			if err != nil {
				return
			}
		} else {
		LOOP:
			err = conn.WriteJSON(Payload{
				Event: "get_ip",
				UUID:  NextUID(),
				Nodes: getIPs(),
				Step:  step,
			})
			if err != nil {
				return
			}

			_, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			r := gjson.GetBytes(p, "..0")
			if r.Get("uuid").String() != nextUUID {
				goto LOOP
			}

			if flow.More {
				nodeJSON := r.Get("node")
				log.Println(nodeJSON, !nodeJSON.Exists(), !nodeJSON.IsArray())
				if !nodeJSON.Exists() || !nodeJSON.IsArray() {
					goto LOOP
				}

				serverNodes := []string{}
				agentNodes := []string{}

				loop := false
				nodeJSON.ForEach(func(key, value gjson.Result) bool {
					node := value.String()
					if !inNodes(node, getIPs()) {
						loop = true
						return false
					}

					if node != masterNode {
						serverNodes = append(serverNodes, node)
					}

					return true
				})

				if loop || len(serverNodes)%2 != 0 {
					goto LOOP
				}

				nodeIPs := getIPs()
				for _, node := range nodeIPs {
					if !inNodes(node, serverNodes) {
						agentNodes = append(agentNodes, node)
					}
				}

				flow.DeployMore(serverNodes, agentNodes)
				fmt.Println(serverNodes, agentNodes)
			} else {
				node := r.Get("node").String()
				if !inNodes(node, getIPs()) {
					goto LOOP
				}

				if flow.Master {
					masterNode = node
				}

				flow.DeployNode(node)
				fmt.Println(node)
			}

			fmt.Println()
		}
	}

	err = conn.WriteJSON(Payload{
		Event: "finish",
		UUID:  NextUID(),
		Step:  11,
	})
	if err != nil {
		return
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		r := gjson.GetBytes(p, "..0")
		if r.Get("uuid").String() != nextUUID {
			continue
		}

		break
	}
}
