package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	src := flag.String("s", "", "原始csv檔案")
	dst := flag.String("d", "", "要複製的csv檔案")
	out := flag.String("o", "default.csv", "輸出檔案")
	flag.Parse()
	if *src == "" || *dst == "" || *out == "" {
		fmt.Println(`
			說明請看 ./copy_csv -h
			完整指令 ./copy_csv -s [原始檔案] -d [複製檔案] -o [輸出檔案]
		`)
		return
	}

	f, err := ioutil.ReadFile(*src)
	if err != nil {
		log.Fatal(err)
	}

	cf, err := ioutil.ReadFile(*dst)
	if err != nil {
		log.Fatal(err)
	}

	// 原始資料
	data := [2][]string{}
	data[0] = []string{}
	data[1] = []string{}

	// 要複製的資料
	cdata := []string{}

	r := csv.NewReader(strings.NewReader(strings.TrimSpace(string(f))))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		data[0] = append(data[0], record[0])
		data[1] = append(data[1], record[1])
	}

	r2 := csv.NewReader(strings.NewReader(strings.TrimSpace(string(cf))))
	for {
		record, err := r2.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		cdata = append(cdata, record[0])
	}

	if len(cdata) == 0 {
		log.Fatal("要複製的csv不能沒有資料！")
	}

	length := -1
	for i, d := range data {
		if len(d) > 0 && d[0] == cdata[0] {
			length = len(cdata)
			data[i] = cdata
		}
		if length > 0 && len(d) != length {
			log.Fatal("要複製的資料量不一樣！")
		} else {
			length = len(d)
		}
	}

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	for i := 0; i < len(data[0]); i++ {
		s := make([]string, 2)
		s[0] = data[0][i]
		s[1] = data[1][i]
		w.Write(s)
		w.Flush()
	}

	fout, err := os.Create(*out)
	defer fout.Close()
	if err != nil {
		fmt.Println(*out, err)
		return
	}
	fout.WriteString(buf.String())
	log.Println("OK")
}
