package main

import (
	"fmt"
	"reflect"
)

// Response 回傳格式
type Response struct {
	Error  *ErrorDetail `json:"error"`
	Result struct {
		ErrorCode string `json:"error_code"`
		ErrorText string `json:"error_text"`
	} `json:"result"`
}

// ErrorDetail 詳細錯誤
type ErrorDetail struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	fmt.Println("Hello World")

	check(struct {
		Error  *ErrorDetail `json:"error"`
		Result struct {
			ErrorCode string `json:"error_code"`
			ErrorText string `json:"error_text"`
			Data      int
		} `json:"result"`
	}{
		Error: nil,
		Result: struct {
			ErrorCode string `json:"error_code"`
			ErrorText string `json:"error_text"`
			Data      int
		}{
			ErrorCode: "12",
			ErrorText: "ok",
			Data:      123,
		},
	})

}

func check(res interface{}) {
	if res == nil {
		return
	}

	shouldResponse := Response{
		Error: nil,
	}
	// var errorCodeName, errorTextName, errorDetailName string

	fmt.Printf("Origin ---> %+v\n", res)
	resStruct := reflect.TypeOf(res)
	resValue := reflect.ValueOf(res)
	num := resStruct.NumField()
	for i := 0; i < num; i++ {
		field := resStruct.Field(i)
		fmt.Printf("Tag %d ---> %+v\n", i, field.Tag.Get("json"))

		switch field.Tag.Get("json") {
		case "error":
		case "result":
			// 用欄位名稱取值
			val := resValue.FieldByName(field.Name)
			fmt.Printf("Value %s ＝＝＝> %v %T\n\n", field.Name, val.Interface(), val.Interface())
			if fmt.Sprintf("%+v", val.Interface()) != "<nil>" {
				subStruct := reflect.TypeOf(val.Interface())
				subValue := reflect.ValueOf(val.Interface())
				subNum := subStruct.NumField()
				for j := 0; j < subNum; j++ {
					subField := subStruct.Field(j)
					switch subField.Tag.Get("json") {
					case "error_code":
						v := subValue.Interface()
						errCode, ok := v.(string)
						if ok {
							shouldResponse.Result.ErrorCode = errCode
						}
					case "error_text":
						v := subValue.Interface()
						errText, ok := v.(string)
						if ok {
							shouldResponse.Result.ErrorText = errText
						}
					}
				}
				fmt.Println(subStruct, subValue, subNum)
			} else {
				fmt.Println("IS NUILL")
			}
		}
	}

	fmt.Printf("\n\nResponse ---> %+v\n", shouldResponse)
}
