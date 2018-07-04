package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Hello World")
	onceCmd := exec.Command("echo", "123")
	keepCmd := exec.Command("ticker")
	onceDone := make(chan error)
	keepDone := make(chan error)
	go runCmd(onceCmd, onceDone)
	go runCmd(keepCmd, keepDone)

	onceErr := <-onceDone
	// onceOutput, err := onceCmd.CombinedOutput()
	// log.Println("once ->", string(onceOutput), err)
	// log.Println("once ->", onceErr, string(onceOutput), err)
	log.Println("once ->", onceErr)

	// keepCmd.Process.Signal(syscall.SIGINT)
	keepErr := <-keepDone
	// keepOutput, err := keepCmd.CombinedOutput()
	// log.Println("keep ->", string(keepOutput), err)
	// log.Println("keep ->", keepErr, string(keepOutput), err)
	log.Println("keep ->", keepErr)

	/// 心得：
	/// 1. Cmd.Output() 或 Cmd.CombinedOutput() 都是直接Run並且取output
	/// 2. 若執行程序非正常結束，Cmd.Output() 不會有輸出
	/// 3. 若執行程序非正常結束，Cmd.CombinedOutput() 還是取得到輸出
	/// 4. Cmd.Run 回傳的error 就可以用來判斷執行狀態，若error包含signal，表示非程式噴出錯誤

}

func runCmd(cmd *exec.Cmd, done chan error) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println("err, ", err)
		if strings.Contains(err.Error(), "signal") {
			close(done)
			return
		}
		done <- err
		return
	}
	close(done)
}
