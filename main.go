package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	// 審視輸入參數
	// log.Println(os.Args)
	// log.Println(os.Args[1:])
	execRoot, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// 取得目標檔案
	targetPath := "."
	if len(os.Args) > 2 {
		targetPath = os.Args[2]
	} else {
		gobuild(execRoot)
		return
	}

	// 檢查是否有家目錄符號
	home := os.Getenv("HOME")
	targetPath = strings.Replace(targetPath, "~", home, 1)
	// log.Println("Root --->", targetPath)

	ticker := time.NewTicker(time.Millisecond * 500)
	rebuild := make(chan error)
	go func() {
		var lastByte int64 = -1
		var lastModTime time.Time
		for range ticker.C {
			// begin := time.Now()
			// 確認.go檔案的佔用大小
			diskByte, lastTime, err := checkGofileByte(targetPath)
			// log.Println(diskByte, err, time.Since(begin))
			if diskByte != lastByte || lastTime.After(lastModTime) {
				rebuild <- err
			}
			lastModTime = lastTime
			lastByte = diskByte
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	var cmd *exec.Cmd

	for {
		select {
		case <-sig:
			ticker.Stop()
			if cmd != nil {
				cmd.Process.Signal(syscall.SIGINT)
				cmd.Wait()
			}
			fmt.Println("\n🐬 ==== 程式結束 ==== 🐬")
			return
		case err := <-rebuild:
			if err == nil {
				if cmd != nil {
					cmd.Process.Signal(syscall.SIGINT)
					cmd.Wait()
				}
				fmt.Println("⚡ ==== 重新編譯==== ⚡")
				cmd = gobuild(execRoot, os.Args[1:]...)
			} else {
				fmt.Println("🔥 ==== 發生錯誤 ==== 🔥", err)
				if cmd != nil {
					cmd.Process.Signal(syscall.SIGINT)
					cmd.Wait()
					cmd = nil
				}
			}
		}
	}
}

func gobuild(execRoot string, args ...string) *exec.Cmd {
	// log.Println("參數, ", args)
	if len(args) > 1 && args[0] == "build" {
		exeFile := execRoot + "/supergo_" + strconv.Itoa(int(time.Now().Unix()))
		// log.Println("執行檔, ", exeFile)
		cmd := exec.Command("go", "build", "-o", exeFile, args[1])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return nil
		}
		if len(args) > 2 {
			cmd = exec.Command(exeFile, args[2:]...)
		} else {
			cmd = exec.Command(exeFile)
		}
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		os.Remove(exeFile)
		return cmd
	}

	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd
}

// checkGofileByte 確認.go檔的總byte大小
func checkGofileByte(root string) (totalByte int64, lastTime time.Time, err error) {
	var f *os.File
	f, err = os.Open(root)
	if err != nil {
		return
	}

	var stats os.FileInfo
	stats, err = f.Stat()
	if err != nil {
		f.Close()
		return
	}
	// log.Println("target path", root)
	if !stats.IsDir() {
		f.Close()
		if strings.HasSuffix(stats.Name(), ".go") {
			return stats.Size(), stats.ModTime(), nil
		}
		return
	}

	var fileList []os.FileInfo
	fileList, err = f.Readdir(-1)
	if err != nil {
		f.Close()
		return
	}
	defer f.Close()

	for i := range fileList {
		file := fileList[i]
		filePath := root + "/" + file.Name()
		// log.Println("check file --->", filePath)
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			// log.Println("target path --->", filePath)
			var extraByte int64
			var extraLastTime time.Time
			extraByte, extraLastTime, err = checkGofileByte(filePath)
			if err != nil {
				log.Println(err)
				continue
			}
			totalByte += extraByte
			if extraLastTime.After(lastTime) {
				lastTime = extraLastTime
			}
		} else if strings.HasSuffix(file.Name(), ".go") {
			// log.Println("get gopher --->", filePath)
			totalByte += file.Size()
			if fileModTime := file.ModTime(); fileModTime.After(lastTime) {
				lastTime = fileModTime
			}
		}
	}
	return
}
