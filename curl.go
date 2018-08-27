package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	for i := 0; i < 1000; i++ {
		go curl()
	}
	time.Sleep(time.Second * 3)
}

func curl() {
	url := "http://10.106.95.97:8000/"

	payload := strings.NewReader("{\n\t\"method\":\"arith.Sum\",\n\t\"params\":{\n\t\t\"A\":0,\n\t\t\"B\":0\n\t}\n}")

	req, _ := http.NewRequest("GET", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "8f08bdc2-c642-102f-7b76-e5d12e5c764f")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
