package main

import (
	"flag"
	"fmt"
	"log"

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

	// 讀取資料
	rawData, err := base.ParseCSVToInstances(*rawDataFile, true)
	if err != nil {
		panic(err)
	}

	// 將資料切分成訓練資料與測試資料
	trainData, testData := base.InstancesTrainTestSplit(rawData, *percent)

	// 建立一個訓練用的工具，cls即是Model
	cls := knn.NewKnnClassifier(*distanceArg, *algorithmArg, *specify)

	// 開始訓練
	err = cls.Fit(trainData)
	if err != nil {
		panic(err)
	}

	// 儲存Model
	err = cls.Save(*outputModelFile)
	if err != nil {
		panic(err)
	}

	log.Println("測試資料 -> ", testData)

	// 預測測試資料
	prediction, err := cls.Predict(testData)
	if err != nil {
		panic(err)
	}
	// 顯示預測結果
	log.Println("預測結果 -> ", prediction)

	// 將測試資料與預測結果，轉化為混淆矩陣
	confusionMat, err := evaluation.GetConfusionMatrix(testData, prediction)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}

	// 顯示評估結果
	fmt.Println(evaluation.GetSummary(confusionMat))
}
