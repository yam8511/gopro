package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var wg = sync.WaitGroup{}

// MessageEndpoint 處理訊息
func MessageEndpoint(bot *tgbotapi.BotAPI, uid int, msg *tgbotapi.Message) {
	if msg.IsCommand() {
		fmt.Printf("[#%d][%s]收到指令, 來自 %s\n指令: %s\n參數: %s\n",
			msg.MessageID,
			msg.Time().Format("2006-01-02 15:04:05-0700"),
			msg.From.String(),
			msg.Command(),
			msg.CommandArguments(),
		)

		rply := NewReplyMessage(msg.Chat.ID, msg.MessageID, "")

		getPwd := func() (string, error) {
			if pwd == "" {
				dir, err := os.Getwd()
				if err != nil {
					return pwd, err
				}

				pwd = dir
			}

			return pwd, nil
		}

		sh := func(raw string) (text string) {
			args := strings.Split(raw, " ")
			var name string
			if len(args) > 0 {
				name = args[0]
				args = args[1:]
			}

			stdout := bytes.NewBuffer(nil)
			stderr := bytes.NewBuffer(nil)
			cm := exec.Command(name, args...)
			cm.Dir = pwd
			cm.Env = os.Environ()
			cm.Stdout = stdout
			cm.Stderr = stderr
			err := cm.Run()
			if err != nil {
				text = fmt.Sprintf(`
						%s指令失敗: %s
						標準輸出:
						%s

						------
						錯誤輸出:
						%s
					`, raw, err.Error(), stdout.String(), stderr.String())
			} else {
				text = fmt.Sprintf(
					"%s指令OK:\n標準輸出:\n%s\n\n------\n錯誤輸出:\n%s",
					raw, stdout.String(), stderr.String(),
				)
			}
			return
		}

		cmdFuncs := map[string]func(){
			"myid": func() {
				rply.Text = fmt.Sprintf("你的ID是(%d)", msg.From.ID)
			},
			"mygroupid": func() {
				if msg.Chat.IsGroup() {
					rply.Text = fmt.Sprintf("你的群組ID是(%d)", msg.Chat.ID)
				} else {
					rply.Text = "目前你所在的聊天室不是群組"
				}
			},
			"ip": func() {
				ips, err := GetLocalIPs()
				if err != nil {
					rply.Text = "本機IP擷取失敗: " + err.Error()
				} else {
					rply.Text = "本機IP如下:\n" + strings.Join(ips, "\n")
				}
			},
			"pwd": func() {
				dir, err := getPwd()
				if err != nil {
					rply.Text = "取當前位子失敗: " + err.Error()
				} else {
					rply.Text = "目前位子是: " + dir
				}
			},
			"ls": func() {
				var err error
				pwd, err = getPwd()
				if err != nil {
					rply.Text = "取當前位子失敗: " + err.Error()
				} else {
					files, err := ioutil.ReadDir(pwd)
					if err != nil {
						rply.Text = "目前位子是:" + pwd + "\n顯示目前資料夾檔案失敗: " + err.Error()
					} else {
						rply.Text = "目前位子是:" + pwd + "\n顯示目前資料夾檔案: "
						for _, f := range files {
							rply.Text += "\n" + f.Name()
							if f.IsDir() {
								rply.Text += "/"
							}
						}
					}
				}
			},
			"cd": func() {
				var err error
				pwd, err = getPwd()
				if err != nil {
					rply.Text = "取當前位子失敗: " + err.Error()
				} else {
					args := msg.CommandArguments()
					switch {
					case strings.HasPrefix(args, "/"):
						pwd = args
					case args == ".":
					case args == ".." && pwd != "/":
						pwd = filepath.Dir(pwd)
					default:
						pwd = filepath.Join(pwd, args)
					}
					rply.Text = "目前位子: " + pwd
				}
			},
			"branch": func() {
				rply.Text = sh("git branch")
			},
			"status": func() {
				rply.Text = sh("git status")
			},
			"reset": func() {
				rply.Text = sh("git reset --hard HEAD")
			},
			"pull": func() {
				var err error
				pwd, err = getPwd()
				if err != nil {
					rply.Text = "取當前位子失敗: " + err.Error()
				} else {
					cmd := exec.Command("cat", pwd+"/.git/HEAD")
					out, err := cmd.Output()
					if err != nil {
						rply.Text = "取當前位子失敗: " + err.Error()
					} else {
						rply.Text += sh("ssh-add -A") + "\n========"
						rply.Text += sh("git pull origin " + strings.TrimSpace(strings.TrimPrefix(string(out), "ref: ")))
					}
				}
			},
			"checkout": func() {
				rply.Text = sh("git checkout " + msg.CommandArguments())
			},
			"build": func() {
				rply.Text = "編譯中，需要稍等一下⛽️"
				bot.Send(rply)

				rply.Text = sh("go build -o program -mod vendor -v .")
			},
			"run": func() {
				stdout := bytes.NewBufferString("")
				stderr := bytes.NewBufferString("")
				fly.Do("run", func() (interface{}, error) {
					if program != nil && (program.ProcessState == nil || !program.ProcessState.Exited()) {
						err := program.Process.Signal(syscall.SIGINT)
						if err != nil {
							rply.Text = "關閉原先伺服器失敗: " + err.Error()
							return nil, err
						}
					}

					program = exec.Command("./program")
					program.Dir = pwd
					program.Env = os.Environ()
					program.Stdout = stdout
					program.Stderr = stderr

					err := program.Start()
					if err != nil {
						rply.Text = "啟動伺服器失敗: " + err.Error()
						return nil, err
					}

					rply.Text = "開始啟動伺服器🚀\n\n"
					bot.Send(rply)
					return nil, err
				})

				program.Wait()
				rply.Text = "伺服器已關閉🛶"
			},
			"stop": func() {
				fly.Do("run", func() (interface{}, error) {
					if program != nil && (program.ProcessState == nil || !program.ProcessState.Exited()) {
						err := program.Process.Signal(syscall.SIGINT)
						if err != nil {
							rply.Text = "關閉原先伺服器失敗: " + err.Error()
							return nil, err
						}
					}
					return nil, nil
				})
			},
			"sh": func() {
				rply.Text = sh(msg.CommandArguments())
			},
		}

		var hintText string = "可輸入以下指令:\n"
		for cmdtxt := range cmdFuncs {
			hintText += "/" + cmdtxt + "\n"
		}
		cmdFuncs["h"] = func() {
			rply.Text = hintText
		}
		cmdFuncs["help"] = func() {
			rply.Text = hintText
		}

		cmd := msg.Command()
		f, ok := cmdFuncs[cmd]
		if ok {
			f()
		} else {
			rply.Text = fmt.Sprintf("未定義指令 [%s], 說明可以輸入 /help", cmd)
		}

		bot.Send(rply)
		return
	}

	fmt.Printf("[#%d][%s]收到新訊息, 來自 %s\n訊息內容: %s\n",
		msg.MessageID,
		msg.Time().Format("2006-01-02 15:04:05-0700"),
		msg.From.String(),
		msg.Text,
	)

	if msg.Sticker != nil {
		sticker := tgbotapi.NewStickerShare(msg.Chat.ID, msg.Sticker.FileID)
		bot.Send(sticker)
		return
	}

	rply := NewReplyMessage(msg.Chat.ID, msg.MessageID, "copy that")
	bot.Send(rply)
}

// CallbackEndpoint 反饋
func CallbackEndpoint(bot *tgbotapi.BotAPI, uid int, cb *tgbotapi.CallbackQuery) {
	fmt.Printf("[#%d][%s]收到新訊息, 來自 %s\n訊息內容: %s\n回傳資料: %s\nChat Instance: %s\nShort Name: %s\ncb ID: %s\ninline message id: %s\n",
		cb.Message.MessageID,
		cb.Message.Time().Format("2006-01-02 15:04:05"),
		cb.From.String(),
		cb.Message.Text,
		cb.Data,
		cb.ChatInstance,
		cb.GameShortName,
		cb.ID,
		cb.InlineMessageID,
	)
	ca := tgbotapi.NewCallback(cb.ID, "跳出CB")
	res, err := bot.AnswerCallbackQuery(ca)
	if err != nil {
		ca = tgbotapi.NewCallbackWithAlert(cb.ID, "CB錯誤: "+err.Error())
		bot.AnswerCallbackQuery(ca)
	} else if !res.Ok {
		fmt.Printf("回傳CB失敗, %#v\n", res)
	}
}

// EditMessageEndpoint 處理訊息
func EditMessageEndpoint(bot *tgbotapi.BotAPI, uid int, msg *tgbotapi.Message) {
	fmt.Printf("[#%d][%s]有新的訊息, 來自 %s\n訊息內容: %s\n",
		msg.MessageID,
		msg.Time().Format("2006-01-02 15:04:05-0700"),
		msg.From.String(),
		msg.Text,
	)
	rply := NewReplyMessage(msg.Chat.ID, msg.MessageID, "catcha that")
	bot.Send(rply)
}

// UpdateMaster 更新入口
func UpdateMaster(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		go func(update tgbotapi.Update) {
			wg.Add(1)
			defer wg.Done()
			defer func(uid int) {
				if r := recover(); r != nil {
					debug.PrintStack()
					log.Printf("[#%d]發生Panic! 原因[%v]\n", uid, r)
				}
			}(update.UpdateID)
			MessageEndpoint(bot, update.UpdateID, update.Message)
		}(update)
	}

	if update.CallbackQuery != nil {
		go func(update tgbotapi.Update) {
			wg.Add(1)
			defer wg.Done()
			defer func(uid int) {
				if r := recover(); r != nil {
					debug.PrintStack()
					log.Printf("[#%d]發生Panic! 原因[%v]\n", uid, r)
				}
			}(update.UpdateID)
			CallbackEndpoint(bot, update.UpdateID, update.CallbackQuery)
		}(update)
	}

	if update.EditedMessage != nil {
		go func(update tgbotapi.Update) {
			wg.Add(1)
			defer wg.Done()
			defer func(uid int) {
				if r := recover(); r != nil {
					log.Printf("[#%d]發生Panic! 原因[%v]\n", uid, r)
				}
			}(update.UpdateID)
			EditMessageEndpoint(bot, update.UpdateID, update.EditedMessage)
		}(update)
	}
}

// WaitShutDown 等待程結束
func WaitShutDown(updateChan tgbotapi.UpdatesChannel) {
	done := make(chan byte)
	go func() {
		updateChan.Clear()
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("執行緒全部完成，結束程序")
	case <-time.After(time.Second * 30):
		log.Println("超過30秒，強制結束程序")
	}
}
