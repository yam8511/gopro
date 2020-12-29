package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

type Payload struct {
	Event string
	Error string
	UUID  string
	Nodes []string
	Step  int
}

type Flow struct {
	Title             string
	More              bool
	NoAction          bool
	Master            bool
	Deploy            func() error
	DeployNode        func(node string) error
	DeployMasterNode  func(node, token string) error
	DeployClusterNode func(nodeMaster, token string, serverNodes, agentNodes []string) error
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
	var token string
	var step int
	NextUID := func() string {
		nextUUID = uuid.New().String()
		return nextUUID
	}

	flows := []Flow{
		{
			Title:      "請選擇一台 Registry Node",
			DeployNode: deployRegistry,
		},
		{
			Title:            "請選擇一台 Master Node",
			DeployMasterNode: deployMaster,
			Master:           true,
		},
		{
			Title:             "請選擇 Server Node (未選作為 Agent Node)",
			More:              true,
			DeployClusterNode: deployCluster,
		},
		{
			Title:      "請選擇一台 Node 作為日誌儲存",
			DeployNode: deployRegistry,
		},
		{
			Title:      "請選擇一台 Node 作為監控儲存",
			DeployNode: deployRegistry,
		},
		{
			Title:      "請選擇一台 Node 作為監控介面(Dashboard)",
			DeployNode: deployRegistry,
		},
		{
			Title:    "部署Database",
			NoAction: true,
			Deploy:   deploy,
		},
		{
			Title:    "部署AI",
			NoAction: true,
			Deploy:   deploy,
		},
		{
			Title:    "部署排程",
			NoAction: true,
			Deploy:   deploy,
		},
		{
			Title:    "部署服務",
			NoAction: true,
			Deploy:   deploy,
		},
	}

	start := func() error {
		return conn.WriteJSON(Payload{
			Event: "start",
			Step:  step,
		})
	}

	end := func() error {
		return conn.WriteJSON(Payload{
			Event: "end",
			Step:  step,
		})
	}

	pushIP := func() error {
		return conn.WriteJSON(Payload{
			Event: "get_ip",
			UUID:  NextUID(),
			Nodes: getIPs(),
			Step:  step,
		})
	}

	pushError := func(err error) {
		log.Println("部署失敗: ", err)
		_ = conn.WriteJSON(Payload{
			Event: "error",
			Error: err.Error(),
			UUID:  NextUID(),
			Step:  step,
		})
	}

	for _, flow := range flows {
		step++

		log.Printf("%d. %s", step, flow.Title)

		if flow.NoAction {
			if start() != nil {
				return
			}

			err = flow.Deploy()
			if err != nil {
				log.Println("部署失敗: ", err)
				pushError(err)
				return
			}

			if end() != nil {
				return
			}
		} else {
		LOOP:
			if pushIP() != nil {
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

				if start() != nil {
					return
				}

				err := flow.DeployClusterNode(masterNode, token, serverNodes, agentNodes)
				if err != nil {
					log.Println("部署失敗: ", err)
					pushError(err)
					return
				}

				if end() != nil {
					return
				}
			} else {
				node := r.Get("node").String()
				if !inNodes(node, getIPs()) {
					goto LOOP
				}

				if start() != nil {
					return
				}

				if flow.Master {
					masterNode = node
					if !isDebug() {
						token = uuid.New().String()
					}
					err = flow.DeployMasterNode(node, token)
				} else {
					err = flow.DeployNode(node)
				}

				if err != nil {
					log.Println("部署失敗: ", err)
					pushError(err)
					return
				}

				if end() != nil {
					return
				}
			}
		}
	}

	step++
	err = conn.WriteJSON(Payload{
		Event: "finish",
		UUID:  NextUID(),
		Step:  step,
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
