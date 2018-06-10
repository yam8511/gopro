package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/sjwhitworth/golearn/base"       // 讀取數據
	"github.com/sjwhitworth/golearn/evaluation" // 評估模型
	"github.com/sjwhitworth/golearn/knn"        // 製作模型
)

func main() {
	rawDataFile := flag.String("r", "", "原始資料的檔案 [.csv]")
	outputModelFile := flag.String("o", "default.model", "輸出訓練模型 [.model]")
	percent := flag.Float64("p", 0.5, "將原始資料的百分之N切分到測試資料")
	specify := flag.Int("s", 2, "特徵點交叉數量")
	distanceArg := flag.String("d", "euclidean", "使用哪一個距離方法 [euclidean, manhattan, cosine]")
	algorithmArg := flag.String("a", "linear", "使用哪一個演算法 [linear, kdtree]")
	flag.Parse()

	if *rawDataFile == "" {
		fmt.Println(`
			請指定輸入資料! ./golearn -r [filename]
			輸入 ./golearn -h ，查看更多詳細資訊
		`)
		return
	}

	allTime := time.Now()

	// 讀取資料
	log.Println("讀取檔案 -> ", *rawDataFile)
	rawData, err := base.ParseCSVToInstances(*rawDataFile, true)
	if err != nil {
		panic(err)
	}

	// 將資料切分成訓練資料與測試資料
	trainData, testData := base.InstancesTrainTestSplit(rawData, *percent)
	_, trainCount := trainData.Size()
	_, testCount := testData.Size()
	log.Printf("切分訓練資料(%d)與測試資料(%d)", trainCount, testCount)

	// 建立一個訓練用的工具，cls即是Model
	cls := knn.NewKnnClassifier(*distanceArg, *algorithmArg, *specify)

	// 開始訓練
	trainTime := time.Now()
	wait := make(chan int)
	done := false
	go func() {
		log.Println("=== 開始訓練 ===")
		err = cls.Fit(trainData)
		if err != nil {
			panic(err)
		}
		wait <- 0
	}()

	for {
		if done {
			break
		}
		select {
		case <-time.After(time.Second * 10):
			log.Println("訓練中...")
		case <-wait:
			done = true
			break
		}
	}
	log.Println("訓練完畢！耗費時間 -> ", time.Since(trainTime))

	// 儲存Model
	log.Println("儲存 Model -> ", *outputModelFile)
	err = cls.Save(*outputModelFile)
	if err != nil {
		panic(err)
	}

	// 預測測試資料
	log.Println("測試資料 -> ", testData)
	log.Println("=== 開始預測 ===")
	var prediction base.FixedDataGrid
	done = false
	go func() {
		prediction, err = cls.Predict(testData)
		if err != nil {
			panic(err)
		}
		wait <- 0
	}()

	for {
		if done {
			break
		}
		select {
		case <-time.After(time.Second * 10):
		case <-wait:
			done = true
			break
		}
	}
	log.Println("預測完畢！耗費時間 -> ", time.Since(trainTime))

	// 顯示預測結果
	log.Println("預測結果 -> ", prediction)

	// 將測試資料與預測結果，轉化為混淆矩陣
	confusionMat, err := evaluation.GetConfusionMatrix(testData, prediction)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}

	// 顯示評估結果
	fmt.Println(evaluation.GetSummary(confusionMat))
	log.Println("程式執行總時間 -> ", time.Since(allTime))
}
