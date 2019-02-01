package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	langJSONFile := flag.String("json", "", "語系的JSON檔案")
	langCSVFile := flag.String("csv", "", "語系的CSV檔案")
	flag.Parse()

	if *langCSVFile == "" || *langJSONFile == "" {
		flag.Usage()
		log.Fatalln("請輸入 csv2json -json [語系.json] -csv [語系.csv]\n例如： ./csv2json -json lang.json -csv lang.csv")
		return
	}

	nowJSLang, err := readLangJSON(*langJSONFile)
	if err != nil {
		log.Fatal(err)
	}

	langCSV, keys, err := readLangCSV(*langCSVFile)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLangCSV(langCSV, keys, nowJSLang)
	if err != nil {
		log.Fatal(err)
	}
}

func readLangJSON(langFile string) (map[string]string, error) {
	rawData, err := ioutil.ReadFile(langFile)
	if err != nil {
		return nil, err
	}

	data := map[string]string{}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, err
	}

	reverseData := map[string]string{}
	for k, v := range data {
		reverseData[v] = k
	}

	return reverseData, nil
}

type lang struct {
	Tw string
	En string
}

func readLangCSV(langFile string) (map[string]lang, []string, error) {
	f, err := os.Open(langFile)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	keys := []string{}
	langDict := map[string]lang{}

	cr := csv.NewReader(f)
	for {
		record, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}

		if len(record) == 0 || record[0] == "" {
			continue
		}

		key := record[0]
		keys = append(keys, key)
		if !strings.HasPrefix(key, "~~~") {
			langObj := lang{Tw: key}
			if len(record) > 1 {
				langObj.En = record[1]
			}
			langDict[key] = langObj
		}
	}

	return langDict, keys, nil
}

func writeLangCSV(
	langCSV map[string]lang,
	keys []string,
	nowJSLang map[string]string,
) error {
	recordTW := []string{}
	recordEN := []string{}

	for i := range keys {
		key := keys[i]
		if strings.HasPrefix(key, "~~~") {
			key = "\n\t" + strings.Replace(key, "~~~", "// ", 1) + "\n"
			recordTW = append(recordTW, key)
			recordEN = append(recordEN, key)
			continue
		}

		trans, ok := langCSV[key]
		if ok {
			jsKey, ok := nowJSLang[key]
			if ok {
				recordTW = append(recordTW, "\t"+jsKey+": '"+trans.Tw+"',")
				en := trans.En
				if en == "" {
					en = trans.Tw
				}
				recordEN = append(recordEN, "\t"+jsKey+": '"+en+"',")
				continue
			}
		}

		recordTW = append(recordTW, "\t'"+key+"': '"+key+"',")
		recordEN = append(recordEN, "\t'"+key+"': '"+key+"',")
	}

	bufTW := []byte(
		"{\n" + strings.Join(recordTW, "\n") + "\n}",
	)

	bufEN := []byte(
		"{\n" + strings.Join(recordEN, "\n") + "\n}",
	)

	var err error
	err = ioutil.WriteFile("./newTW.json", bufTW, os.FileMode(0666))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./newEN.json", bufEN, os.FileMode(0666))
	if err != nil {
		return err
	}

	return nil
}
