package main

import (
	"fmt"
	"log"

	"github.com/sjwhitworth/golearn/base"       // 讀取數據
	"github.com/sjwhitworth/golearn/evaluation" // 評估模型
	"github.com/sjwhitworth/golearn/knn"        // 製作模型
)

func main() {
	// 讀取資料
	rawData, err := base.ParseCSVToInstances("dataset/abalone/data.csv", true)
	if err != nil {
		panic(err)
	}

	// 將資料切分成訓練資料與測試資料
	trainData, testData := base.InstancesTrainTestSplit(rawData, 0.1)

	// 建立一個訓練用的工具，cls即是Model
	cls := knn.NewKnnClassifier("euclidean", "linear", 1)

	// 開始訓練
	err = cls.Fit(trainData)
	if err != nil {
		panic(err)
	}

	// 儲存Model
	cls.Save("abalone.model")

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
