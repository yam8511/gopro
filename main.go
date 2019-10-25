package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	rds := convert("./cash_record.csv")
	txs := convert("./transaction.csv")

	var mainTxs, subTxs []string
	lenTXS := len(txs)
	lenRDS := len(rds)
	if lenRDS < lenTXS {
		mainTxs = make([]string, lenTXS)
		subTxs = make([]string, lenRDS)
		copy(mainTxs, txs)
		copy(subTxs, rds)
	} else {
		mainTxs = make([]string, lenRDS)
		subTxs = make([]string, lenTXS)
		copy(mainTxs, rds)
		copy(subTxs, txs)
	}

	subMap := map[string]int{}
	for _, tx := range subTxs {
		subMap[tx] = 0
	}

	extraTx := []string{}
	for _, tx := range mainTxs {
		_, ok := subMap[tx]
		if !ok {
			extraTx = append(extraTx, tx)
		}
	}

	log.Println(extraTx)
	log.Println(len(extraTx))

	err := ioutil.WriteFile("extra.txt", []byte(strings.Join(extraTx, "\n")), 0777)
	if err != nil {
		log.Fatal("寫入檔案Error", err)
	}
}

func convert(name string) []string {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	bs := strings.Split(string(b), "\n")
	lenBS := len(bs)
	if bs[len(bs)-1] == "" {
		lenBS = len(bs) - 1
	}
	bbs := make([]string, lenBS)
	copy(bbs, bs)
	return bbs
}
