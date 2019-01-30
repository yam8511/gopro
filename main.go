package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func main() {
	goodData, _, err := readCsvData("./good.csv")
	if err != nil {
		log.Println("Load Good Data", err)
		return
	}

	langData, keys, err := readCsvData("./lang.csv")
	if err != nil {
		log.Println("Load Lang Data", err)
		return
	}

	log.Println(len(langData), len(goodData))
	for key := range langData {
		val, ok := goodData[key]
		if ok {
			langData[key] = val
		}
	}

	w, err := os.OpenFile("./zzz.csv", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Println("Open Err ", err)
		return
	}
	defer w.Close()

	cw := csv.NewWriter(w)
	for _, key := range keys {
		val, ok := langData[key]
		if !ok {
			continue
		}
		err = cw.Write([]string{key, val})
		if err != nil {
			log.Println("Write Err ", err)
			return
		}
	}
	cw.Flush()
}

func readCsvData(filename string) (dict map[string]string, keys []string, e error) {
	r, err := os.Open(filename)
	if err != nil {
		log.Println("Open Err ", err)
		e = err
		return
	}
	defer r.Close()

	cr := csv.NewReader(r)
	data, err := cr.ReadAll()
	if err != nil {
		log.Println("Read Err ", err)
		e = err
		return
	}

	dict = map[string]string{}
	keys = []string{}
	for _, d := range data {
		if len(d) > 0 {
			var key, val string
			key = d[0]
			val = strings.Join(d[1:], ",")

			_, ok := dict[key]
			if !ok {
				keys = append(keys, key)
			}
			dict[key] = val
		}
	}

	return
}
