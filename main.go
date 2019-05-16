package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
)

func main() {
	startup()
}

func startup() {
	env := os.Environ()

	isNormal := true
	for i := range env {
		if env[i] == "ZZ=1" {
			isNormal = false
			break
		}
	}

	if !isNormal {
		hackMode()
		return
	}

	normalMode()
}

func normalMode() {
	if len(os.Args) == 0 {
		return
	}

	b := os.Args[0]
	hackArgs := []string{b}
	pid, err := syscall.ForkExec(b, hackArgs, &syscall.ProcAttr{
		Env: append(os.Environ(), "ZZ=1"),
	})

	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Pid: ", pid)
}

func hackMode() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/hack", func(w http.ResponseWriter, r *http.Request) {
		info := "Args : " + strings.Join(os.Args, " ") + "\n"
		info += "Env : " + strings.Join(os.Environ(), "\n") + "\n"
		w.Write([]byte(info))
	})

	http.ListenAndServe(":8000", nil)
}
