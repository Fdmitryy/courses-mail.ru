package main

// сюда писать код

import (
	"context"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "803773001:AAFgW1dinrfRbckpox5YYlxCC3KBSVRprMs"

	// урл выдаст вам игрок или хероку
	WebhookURL = "https://taskerhwbot.herokuapp.com/"
)

var commands = map[string]func(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	"/tasks":    getTasks,
	"/new":      newTask,
	"/assign":   assign,
	"/unassign": unassign,
	"/resolve":  resolve,
	"/my":       getMyTasks,
	"/owner":    getTasksByMe,
}

type userTask struct {
	ChatId   int64
	Task     string
	UserName string
}

var userTasks []userTask
var allTasks = map[int64][]string{}
var id int64 = 0

func startTaskBot(ctx context.Context) error {
	// сюда пишите ваш код
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return err
	}
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		return err
	}
	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":8081", nil)
	fmt.Println("start listen :8081")
	for update := range updates {
		text := update.Message.Text
		str := strings.Split(text, " ")
		comm := str[0]
		if strings.Contains(comm, "_") {
			newcomm := strings.Split(text, "_")
			comm = newcomm[0]
		}
		if _, exist := commands[comm]; !exist {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "unknown command\n"))
			continue
		}
		commands[comm](update, bot)
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		fmt.Println("something bad happened")
	}
}

func getTasks(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatId := update.Message.Chat.ID
	var result string
	var res2 string
	var count int
	var sortedByTaskId = make(map[int]string)
	for _, val := range allTasks {
		for _, value := range val {
			str := strings.Split(value, ".")
			taskId, _ := strconv.Atoi(str[0])
			sortedByTaskId[taskId] = value
		}
	}
	sorted := make([]int, 0, len(sortedByTaskId))
	for tId := range sortedByTaskId {
		sorted = append(sorted, tId)
	}
	sort.Ints(sorted)
	for _, s := range sorted {
		v := sortedByTaskId[s]
		str := strings.Split(v, ".")
		taskId := str[0]
		var flag bool
		for _, val := range userTasks {
			if strings.Contains(val.Task, taskId+".") {
				if val.ChatId == chatId {
					if count == 0 {
						res2 = "assignee: я\n" + "/unassign_" + taskId + " /resolve_" + taskId
					} else {
						res2 = "\nassignee: я\n" + "/unassign_" + taskId + " /resolve_" + taskId + "\n"
					}
				} else {
					res2 = "assignee: @" + val.UserName
				}
				flag = true
			}
		}
		if flag {
			if count == 0 {
				result += v + res2
			} else {
				result += "\n\n" + v + res2
			}
		} else {
			if count == 0 {
				result += v + "/assign_" + taskId
			} else {
				result += "\n\n" + v + "/assign_" + taskId
			}

		}
		count++
	}
	if result == "" {
		bot.Send(tgbotapi.NewMessage(chatId, "Нет задач"))
	} else {
		bot.Send(tgbotapi.NewMessage(chatId, result))
	}
}

func newTask(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	text := update.Message.Text
	chatId := update.Message.Chat.ID
	id++
	str := strings.Split(text, " ")
	args := strings.Join(str[1:], " ")
	allTasks[chatId] = append(allTasks[chatId], strconv.Itoa(int(id))+". "+args+" by @"+update.Message.Chat.UserName+"\n")
	bot.Send(tgbotapi.NewMessage(chatId, `Задача "`+args+`" создана, id=`+strconv.Itoa(int(id))))
}

func assign(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	text := update.Message.Text
	chatId := update.Message.Chat.ID
	str := strings.Split(text, "_")
	taskId := str[1]
	var creatorId int64
	var task string
	for k, value := range allTasks {
		for _, val := range value {
			if strings.Contains(val, taskId+".") {
				task = val
				creatorId = k
				break
			}
		}
	}
	for k, arr := range userTasks {
		if strings.Contains(arr.Task, taskId+".") {
			creatorId = arr.ChatId
			userTasks = append(userTasks[:k], userTasks[k+1:]...)
			break
		}
	}
	userTasks = append(userTasks, userTask{ChatId: chatId, UserName: update.Message.Chat.UserName, Task: task})
	mas := strings.Split(task, " ")
	mas = mas[1 : len(mas)-2]
	res1 := strings.Join(mas, " ")
	if creatorId != chatId {
		bot.Send(tgbotapi.NewMessage(creatorId, "Задача "+`"`+res1+`" `+"назначена на @"+update.Message.Chat.UserName))
	}
	bot.Send(tgbotapi.NewMessage(chatId, "Задача "+`"`+res1+`" `+"назначена на вас"))

}

func unassign(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	text := update.Message.Text
	chatId := update.Message.Chat.ID
	str := strings.Split(text, "_")
	taskId := str[1]
	var creatorId int64
	for k, arr := range userTasks {
		if strings.Contains(arr.Task, taskId+".") && arr.ChatId == chatId {
			userTasks = append(userTasks[:k], userTasks[k+1:]...)
			for key, value := range allTasks {
				for _, val := range value {
					if val == arr.Task {
						creatorId = key
						break
					}
				}
			}
			mas := strings.Split(arr.Task, " ")
			mas = mas[1 : len(mas)-2]
			res1 := strings.Join(mas, " ")
			bot.Send(tgbotapi.NewMessage(creatorId, `Задача `+`"`+res1+`"`+` осталась без исполнителя`))
			bot.Send(tgbotapi.NewMessage(chatId, "Принято"))
			return
		}
	}
	bot.Send(tgbotapi.NewMessage(chatId, `Задача не на вас`))

}

func resolve(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	text := update.Message.Text
	chatId := update.Message.Chat.ID
	str := strings.Split(text, "_")
	taskId := str[1]
	var who string
	var creatorId int64
	for k, arr := range userTasks {
		if strings.Contains(arr.Task, taskId+".") && arr.ChatId == chatId {
			who = arr.UserName
			userTasks = append(userTasks[:k], userTasks[k+1:]...)
			for key, value := range allTasks {
				for _, val := range value {
					if val == arr.Task {
						delete(allTasks, key)
						creatorId = key
						break
					}
				}
			}
			mas := strings.Split(arr.Task, " ")
			mas = mas[1 : len(mas)-2]
			res1 := strings.Join(mas, " ")
			if creatorId != chatId {
				bot.Send(tgbotapi.NewMessage(creatorId, `Задача `+`"`+res1+`"`+` выполнена @`+who))
			}
			bot.Send(tgbotapi.NewMessage(chatId, `Задача `+`"`+res1+`"`+` выполнена`))
			return
		}
	}
}

func getMyTasks(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var result string
	var count int
	chatId := update.Message.Chat.ID
	for _, val := range userTasks {
		if val.ChatId == chatId {
			str := strings.Split(val.Task, ".")
			taskId := str[0]
			if count == 0 {
				result += val.Task + "/unassign_" + taskId + " /resolve_" + taskId
			} else {
				result += "\n" + val.Task + "/unassign_" + taskId + " /resolve_" + taskId
			}
			count++
		}
	}
	if result == "" {
		bot.Send(tgbotapi.NewMessage(chatId, "Нет задач"))
	} else {
		bot.Send(tgbotapi.NewMessage(chatId, result))
	}
}

func getTasksByMe(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatId := update.Message.Chat.ID
	var result string
	var count int
	for key, value := range allTasks {
		for _, val := range value {
			if key == chatId {
				str := strings.Split(val, ".")
				taskId := str[0]
				if count == 0 {
					result += val + "/assign_" + taskId
				} else {
					result += "\n" + val + "/assign_" + taskId
				}
				count++
			}
		}
	}
	if result == "" {
		bot.Send(tgbotapi.NewMessage(chatId, "Нет задач"))
	} else {
		bot.Send(tgbotapi.NewMessage(chatId, result))
	}
}
