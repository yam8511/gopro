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
	// var i interface{}
	// fmt.Printf("%+v\n", reflect.ValueOf(i).IsValid())
	// return

	check(struct {
		Error  *ErrorDetail `json:"error"`
		Result struct {
			ErrorCode string `json:"error_code"`
			ErrorText string `json:"error_text"`
			Data      int
		} `json:"result"`
	}{
		Error: &ErrorDetail{
			Code:    500,
			Message: "123",
			Data:    "OKOKK",
		},
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

		val := resValue.FieldByName(field.Name)
		v := val.Interface()
		fmt.Printf("Value %s ＝＝＝> %v %T\n", field.Name, v, v)
		if reflect.ValueOf(v).IsValid() {
			fmt.Printf("Value %s ＝＝＝> valid: %v\n", field.Name, reflect.ValueOf(v).IsValid())
			if reflect.ValueOf(v).Kind() == reflect.Ptr {
				fmt.Printf("Value %s ＝＝＝> nil: %v\n\n", field.Name, reflect.ValueOf(v).IsNil())
				if reflect.ValueOf(v).IsNil() {
					continue
				}
			}
		} else {
			fmt.Printf("Value %s ＝＝＝> valid: %v\n", field.Name, reflect.ValueOf(v).IsValid())
			continue
		}

		switch field.Tag.Get("json") {
		case "error":
			sv, ok := v.(*ErrorDetail)
			if ok {
				shouldResponse.Error = sv
			}
		case "result":
			// 用欄位名稱取值
			if fmt.Sprintf("%+v", v) != "<nil>" {
				subStruct := reflect.TypeOf(v)
				subValue := reflect.ValueOf(v)
				subNum := subStruct.NumField()
				for j := 0; j < subNum; j++ {
					subField := subStruct.Field(j)
					switch subField.Tag.Get("json") {
					case "error_code":
						sf := subValue.FieldByName(subField.Name)
						sv := sf.Interface()
						errCode, ok := sv.(string)
						fmt.Println("ErrorCode --->", sv)
						if ok {
							shouldResponse.Result.ErrorCode = errCode
						}
					case "error_text":
						sf := subValue.FieldByName(subField.Name)
						sv := sf.Interface()
						errText, ok := sv.(string)
						fmt.Println("ErrorText --->", sv)
						if ok {
							shouldResponse.Result.ErrorText = errText
						}
					}
				}
				fmt.Println(subStruct, subValue, subNum)
			}
		}
	}

	fmt.Printf("\n\nResponse ---> %+v\n", shouldResponse)
	if shouldResponse.Error != nil {
		fmt.Println("Response ---> ", shouldResponse.Error.Code, shouldResponse.Error.Message, shouldResponse.Error.Data)
	}
}
