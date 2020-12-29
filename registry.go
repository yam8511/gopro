package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/txn2/txeh"
)

func deployRegistry(node string) error {
	hasLocal := false
	nodes := getIPs()
	errCh := make(chan error)
	for _, ip := range nodes {
		go func(ip string) {
			if isLocalIP(ip) {
				if hasLocal {
					errCh <- nil
					return
				}
				ip = "127.0.0.1"
				node = "127.0.0.1"
				hasLocal = true
			}

			err := curlRegistry(ip, node)
			errCh <- err
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

func curlRegistry(ip, node string) error {
	link := url.URL{
		Scheme: "http",
		Host:   ip + PORT,
		Path:   "/deploy/registry",
	}
	query := link.Query()
	query.Set("node", node)
	link.RawQuery = query.Encode()

	method := http.MethodPost
	client := &http.Client{Timeout: time.Minute}
	req, err := http.NewRequest(method, link.String(), nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http not 200: %s", res.Status)
	}

	txt := res.Header.Get("error")
	if txt != "" {
		return fmt.Errorf(txt)
	}

	return nil
}

func processRegistry(c *gin.Context) {
	node := c.Query("node")
	if node == "" {
		c.Header("error", "No IP")
		return
	}

	if net.ParseIP(node) == nil {
		c.Header("error", "IP Invalid")
		return
	}

	ho, err := txeh.NewHostsDefault()
	if err != nil {
		c.Header("error", err.Error())
		return
	}
	ho.WriteFilePath = "./host"

	ok, ip, _ := ho.HostAddressLookup("registry")
	if ok && ip == node {
		return
	}

	ho.AddHost(node, "registry")
	err = ho.Save()
	if err != nil {
		c.Header("error", err.Error())
	}
}
