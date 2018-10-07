package main

import (
	"fmt"
	"text/template"
)

func main() {
	fmt.Println("Hello World")

	i18n := NewI18n()
	// locale := map[string]map[string]string{
	// 	"en": map[string]string{
	// 		"hello":   "Hello {{ .Name }}",
	// 		"greeter": "Nice to meet you",
	// 	},
	// 	"cn": map[string]string{
	// 		"hello":   "哈囉 {{ .Name }}",
	// 		"greeter": "初次見面",
	// 	},
	// }
	// i18n.Load(locale)
	err := i18n.LoadJSON("./assets/translate.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]string{"Name": "Zuolar"}
	fmt.Println("顯示文字 ---> ", i18n.Trans("hello", data))
	fmt.Println("顯示文字 ---> ", i18n.Trans("greeter"))

	i18n.SetLocale("cn")
	fmt.Println("顯示文字 ---> ", i18n.Trans("hello", data))
	fmt.Println("顯示文字 ---> ", i18n.Trans("greeter"))
}

// NewI18n 建立一個翻譯機
func NewI18n() *I18n {
	return &I18n{
		tmp:          template.New("i18n"),
		fallbackLang: "en",
		currentLang:  "en",
	}
}
