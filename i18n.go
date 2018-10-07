package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type translation struct {
	Name     string
	Template string
}

type translationSet struct {
	Language   string
	Dictionary []translation
}

// I18n 翻譯機
type I18n struct {
	tmp *template.Template
	// dictionary   []translationSet
	fallbackLang string
	currentLang  string
}

// LoadJSON 載入字典檔
func (i18n *I18n) LoadJSON(jsonFile string) error {
	f, err := os.Open(jsonFile)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var dicts []translationSet
	err = json.Unmarshal(jsonData, &dicts)
	if err != nil {
		var dict translationSet
		err = json.Unmarshal(jsonData, &dict)
		if err != nil {
			return err
		}
		dicts = append(dicts, dict)
	}

	i18n.load(dicts)
	return nil
}

// Load 載入字典檔
func (i18n *I18n) Load(data map[string]map[string]string) {
	if data != nil {
		dicts := []translationSet{}
		for lang, dictionary := range data {
			dict := translationSet{
				Language: lang,
			}
			for name, tpl := range dictionary {
				dict.Dictionary = append(dict.Dictionary, translation{
					Name:     name,
					Template: tpl,
				})
				// i18n.tmp.New(lang + "." + name).Parse(tpl)
			}
		}
		i18n.load(dicts)
	}
}

func (i18n *I18n) load(dicts []translationSet) {
	for i := range dicts {
		dict := &dicts[i]
		for j := range dict.Dictionary {
			translation := &dict.Dictionary[j]
			i18n.tmp.New(dict.Language + "." + translation.Name).Parse(translation.Template)
		}
	}
}

// Trans 翻譯
func (i18n *I18n) Trans(tpl string, data ...interface{}) string {
	t := i18n.tmp.Lookup(i18n.currentLang + "." + tpl)
	if t == nil {
		t = i18n.tmp.Lookup(i18n.fallbackLang + "." + tpl)
		if t == nil {
			return tpl
		}
	}
	b := new(strings.Builder)
	if len(data) > 0 {
		t.Execute(b, data[0])
	} else {
		t.Execute(b, nil)
	}
	return b.String()
}

// SetLocale 設定目前語系
func (i18n *I18n) SetLocale(lang string) {
	i18n.currentLang = lang
}

// SetFallbackLocale 設定無效之後的語系
func (i18n *I18n) SetFallbackLocale(lang string) {
	i18n.fallbackLang = lang
}
