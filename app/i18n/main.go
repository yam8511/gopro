package i18n

import (
	"bytes"
	"strings"
	"text/template"

	lua "github.com/yuin/gopher-lua"
)

const i18nRoot = "app/i18n"
const luaFile = "main.lua"

// Trans 翻譯
func Trans(lang, key string, args map[string]string) string {
	split := func(l *lua.LState) int {
		key := l.GetGlobal("key").String()
		keys := strings.Split(key, ".")

		// === Reverse ===
		last := len(keys) - 1
		for i := 0; i < len(keys)/2; i++ {
			keys[i], keys[last-i] = keys[last-i], keys[i]
		}

		table := &lua.LTable{}
		for i := range keys {
			table.RawSet(lua.LNumber(i), lua.LString(keys[i]))
		}

		// l.Push = return
		l.Push(table)
		return 1 // 指定回傳幾個資料
	}

	mapping := func(l *lua.LState) int {
		output := l.GetGlobal("output")
		if output.Type() != lua.LTString {
			l.Push(output)
		} else {
			out := output.String()
			tpl, err := template.New("/").Parse(out)
			if err == nil {
				w := bytes.NewBuffer([]byte{})
				err = tpl.Execute(w, args)
				if err == nil {
					out = w.String()
				}
			}

			out = strings.ReplaceAll(out, "<no value>", "")
			output = lua.LString(out)
			l.Push(output)
		}

		return 1
	}

	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("root", lua.LString(i18nRoot))
	L.SetGlobal("lang", lua.LString(lang))
	L.SetGlobal("key", lua.LString(key))
	L.SetGlobal("split", L.NewFunction(split))
	L.SetGlobal("mapping", L.NewFunction(mapping))

	if err := L.DoFile(i18nRoot + "/" + luaFile); err != nil {
		panic(err)
	}

	return lua.LVAsString(L.GetGlobal("output"))
}
