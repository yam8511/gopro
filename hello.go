package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorhill/cronexpr"
)

// Job 任務
type Job struct {
	Name           string
	Cron           string
	Command        string
	Args           []string
	Env            []string
	IsKeep         bool
	IsOverlapping  bool
	Enable         bool
	cmd            map[int64]*exec.Cmd
	cronExpression *cronexpr.Expression
	cancel         chan int
	mx             *sync.RWMutex
}

var globalJobs map[string]*Job
var globalSignal os.Signal

func main() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	iniJobs := map[string]*Job{
		"job echo": &Job{
			Name:    "job echo",
			Cron:    "* * * * s* *",
			Command: "echo",
			Args:    []string{"10"},
			Enable:  true,
		},
		// &Job{
		// 	Name:    "job ticker",
		// 	Cron:    "*/5 * * * *aa*",
		// 	Command: "ticker",
		// 	Args:    []string{"-s=2"},
		// 	Enable:  true,
		// 	IsKeep:  true,
		// },
	}

	// 四秒後，刷新排程
	go func() {
		<-time.After(time.Second * 4)
		refreshSchedule()
	}()

	// 排程初始化
	err := ScheduleInit(iniJobs)
	if err != nil {
		log.Fatal("init -> ", err)
	}

	globalJobs = iniJobs

	// 執行背景
	for _, job := range globalJobs {
		go RunJob(job)
	}

	// 等待中斷訊號
	globalSignal = <-sig
	log.Println("Getcha Signal ->", globalSignal)

	// 結束背景
	for _, job := range globalJobs {
		// 等待背景結束
		waitJobDone(job)
	}

	log.Println("Exit")
}

// ScheduleInit 跑背景
func ScheduleInit(jobs map[string]*Job) (err error) {
	reg := regexp.MustCompile(`([a-zA-Z]+|\s{2,})`)

	for _, job := range jobs {
		job.cmd = map[int64]*exec.Cmd{}
		job.cancel = make(chan int)
		job.mx = new(sync.RWMutex)

		spec := job.Cron
		// 把多餘字母轉成空格
		spec = string(reg.ReplaceAll([]byte(spec), []byte(" ")))
		// 把多餘空格轉成一格空格
		spec = string(reg.ReplaceAll([]byte(spec), []byte(" ")))
		// 去除前後空白
		spec = strings.TrimSpace(spec)

		// 如果星星數量是六顆，依照原本規則是會有秒數，但是這套件需要七顆才有秒數
		stars := strings.Split(spec, " ")
		if len(stars) == 6 {
			spec += " *"
		}
		job.cronExpression, err = cronexpr.Parse(spec)
		if err != nil {
			return
		}
		job.Cron = spec
		log.Println(job.Name, job.Cron)
	}
	return
}

// RunJob 跑背景
func RunJob(job *Job) {
	for {
		// 現在時間
		now := time.Now()
		// 下次執行時間
		next := job.cronExpression.Next(now)

		select {
		case <-time.After(next.Sub(now)):
			// 如果沒啟動，則繼續等下一次
			if !job.Enable {
				continue
			}

			// 等待執行時間到，先檢查是否還在執行中
			var doing bool
			job.mx.RLock()
			doing = len(job.cmd) > 0
			job.mx.RUnlock()

			// 如果是不能重複執行的，而且在執行則等下一次
			if !job.IsOverlapping && doing {
				continue
			}

			// 如果是常駐的工作，則如果還在執行也不能執行
			// 但是需要紀錄時間
			if job.IsKeep && doing {
				// 紀錄執行時間
				record()
				continue
			}

			// 建立指令
			cmd := createCmd(job)

			// 儲存到 job.cmd 中
			job.mx.Lock()
			job.cmd[next.Unix()] = cmd
			job.mx.Unlock()

			go func(cmd *exec.Cmd, excuteTime time.Time) {
				done := make(chan error)

				// 紀錄執行時間
				record()

				// 開始執行
				go func() {
					err := cmd.Run()
					if err != nil {
						if strings.HasPrefix(err.Error(), "signal") {
							done <- nil
						} else {
							done <- err
						}
						return
					}
					done <- err
				}()

				// 等待執行完成
				<-done

				// 紀錄結束時間與狀態
				record()

				// 關閉通道
				close(done)

				// 從 job.cmd 中移除
				job.mx.Lock()
				delete(job.cmd, next.Unix())
				job.mx.Unlock()
			}(cmd, next)
		case <-job.cancel:
			// 如果收到取消的訊號，則結束
			log.Printf("[%s] 收到結束訊號，kill 程序\n", job.Name)
			job.mx.Lock()
			for _, cmd := range job.cmd {
				cmd.Process.Signal(syscall.SIGINT)
			}
			job.mx.Unlock()
			return
		}
	}
}

func createCmd(job *Job) (cmd *exec.Cmd) {
	cmd = exec.Command(job.Command, job.Args...)
	cmd.Env = append(os.Environ(), job.Env...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return
}

func record() {
}

// 等待背景結束
func waitJobDone(job *Job) {
	close(job.cancel)
	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		job.mx.RLock()
		count := len(job.cmd)
		job.mx.RUnlock()
		if count == 0 {
			ticker.Stop()
			return
		}
	}
	return
}

// 刷新排程背景
func refreshSchedule() {
	if globalSignal != nil {
		return
	}
	newJobs := map[string]*Job{
		"job echo": &Job{
			Name:    "job echo",
			Cron:    "*/5 * * * s* *",
			Command: "sleep",
			Args:    []string{"10"},
			Enable:  true,
		},
		"job ticker": &Job{
			Name:    "job ticker",
			Cron:    "*/5 * * * *aa*",
			Command: "ticker",
			Args:    []string{"-s=2"},
			Enable:  true,
			IsKeep:  true,
		},
	}

	ScheduleInit(newJobs)
	for _, newJob := range newJobs {
		if job, ok := globalJobs[newJob.Name]; ok {
			if job.Cron != newJob.Cron {
				log.Printf(">>>>> 因為 Cron 時間不一樣，所以重起背景。 等待舊的 [%s] 結束 <<<<<<\n", job.Name)
				waitJobDone(job)
			}
		}
		log.Printf(">>>>> 啟動新背景 [%s]  <<<<<<\n", newJob.Name)
		go RunJob(newJob)
	}
	globalJobs = newJobs
}
