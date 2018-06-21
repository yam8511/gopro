package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	input := flag.String("f", "", "題庫檔")
	flag.Parse()

	// 若輸入檔案為空字串
	if *input == "" {
		fmt.Printf(`
	* 請提供檔案, -f 指定檔案, -h 詳見指令 *

	Version 1.0.1

`)
		return
	}

	b, err := ioutil.ReadFile(*input)
	if err != nil {
		log.Fatal("讀檔錯誤 ->", *input)
	}

	var (
		content   = string(b)
		all       = strings.Split(strings.TrimSpace(content), "\n")
		questions = [][]string{}
		question  []string
	)

	// 切題目
	for _, line := range all {
		// 處理前後空白
		trimLine := strings.TrimSpace(line)
		if trimLine == "" {
			continue
		}

		// 如果是Q開頭，但不是QUESTION，表示新的題目
		if strings.HasPrefix(trimLine, "Q") && !strings.HasPrefix(trimLine, "QUESTION") {
			if len(question) > 0 {
				questions = append(questions, question)
			}
			question = []string{line}
			continue
		}

		// 句子一律往裡面塞
		question = append(question, line)
	}

	if len(question) > 0 {
		questions = append(questions, question)
	}

	var (
		files []string
	)
	// 開始寫入文字檔
	for _, question := range questions {
		q := question[0]
		question = question[1:]
		c := []byte(strings.Join(question, "\n"))
		prefix := strings.TrimSuffix(*input, ".txt")
		filename := prefix + "_" + q + ".txt"
		err := ioutil.WriteFile(filename, c, 0644)
		if err != nil {
			log.Printf("%s 寫入檔案錯誤！ (%s)\n", question[0], err.Error())
		} else {
			files = append(files, filename)
		}
	}

	log.Println("========== 輸出以下題目檔 =========")
	for _, filename := range files {
		log.Println(filename)
	}
	log.Printf("切題共 %d 題, 輸出共 %d 題\n", len(questions), len(files))
}
