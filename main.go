package main

import (
	"bytes"
	"log"
	"text/template"
)

func main() {
	// 假資料
	// jsonStr := `{"User":"Bot", "GameID":1}`
	// type Record struct {
	// 	User   string `json:"User"`
	// 	GameID int    `json:"GameID"`
	// }

	// data := Record{}
	// err := json.Unmarshal([]byte(jsonStr), &data)
	// if err != nil {
	// 	log.Fatal("Decode", err)
	// }
	data := map[string]interface{}{
		"User":   "Zuolar",
		"GameID": "gameID",
	}
	var err error

	def := map[string]string{
		"tpl1": `{{ .User | df "Nobody" }} 新增遊戲 {{ .GameID | df 0 }}`,
		"tpl2": `由 {{ .Username | df "Nobody" }} 刪除遊戲 {{ .GameIDs | df 0 }}`,
	}

	rootTpl := template.New("").Funcs(map[string]interface{}{
		"df": func(fallback interface{}, source interface{}) interface{} {
			if source == nil {
				return fallback
			}
			return source
		},
	})
	for key, val := range def {
		_, err = rootTpl.New(key).Parse(val)
		if err != nil {
			log.Fatal("T1 Parse", err)
		}
	}

	w := new(bytes.Buffer)

	w.Reset()
	err = rootTpl.ExecuteTemplate(w, "tpl1", data)
	if err != nil {
		log.Fatal("Exec", err)
	}
	log.Println(w.String())

	w.Reset()
	err = rootTpl.ExecuteTemplate(w, "tpl2", data)
	if err != nil {
		log.Fatal("Exec", err)
	}
	log.Println(w.String())

}
