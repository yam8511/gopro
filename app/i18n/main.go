package i18n

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

const i18nRoot = "app/i18n"
const luaFile = "main.lua"

var proto *lua.FunctionProto

// 提前編譯
func init() {
	source, err := ioutil.ReadFile(i18nRoot + "/" + luaFile)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewBuffer(source)

	chunk, err := parse.Parse(reader, reader.String())
	if err != nil {
		panic(err)
	}

	proto, err = lua.Compile(chunk, reader.String())
	if err != nil {
		panic(err)
	}
}

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

	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	err := L.PCall(0, lua.MultRet, nil)
	if err != nil {
		panic(err)
	}

	// if err := L.DoFile(i18nRoot + "/" + luaFile); err != nil {
	// 	panic(err)
	// }

	return lua.LVAsString(L.GetGlobal("output"))
}
