package main

import (
	"log"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var botUser *api.BotAPI
var questionsData map[string][]string = map[string][]string{
	"Сколько тебе лет?": []string{
		"7-8",
		"9-12",
		"12-15",
		"16-18",
	},
	"Что ты предпочтешь?": []string{
		"Покормить улитку",
		"Пройти полезный бизнес-трейнинг",
		"Наругать губку",
		"Украсть рецепт бургера",
	},
	"Что бы ты сьел на обед?": []string{
		"Красбургер",
		"Я сьем все!",
		"Кравиоли",
		"Мясо в ведре",
	},
}
var questions []string = []string{
	"Сколько тебе лет?",
	"Что ты предпочтешь?",
	"Что бы ты сьел на обед?",
}
var users map[int64]*user
var characters []string = []string{"Губка Боб", "Патрик", "Сквидварт", "Планктон"}

func main() {
	bot, err := api.NewBotAPI("TOKEN")

	if err != nil {
		log.Panic(err)
	}

	botUser = bot

	botUser.Debug = true

	log.Printf("Authorized on account %s", botUser.Self.UserName)
	
	users = make(map[int64]*user)

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go onNewMessage(update.Message)
		}

		if update.CallbackQuery != nil {
			callback := api.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println(err)
			}
			go onInlineData(update.CallbackQuery)
		}
	}
}

func PrintResult(id int64) {
	user := users[id]
	resultName := "spongebob"
	for characterName := range user.Result {
		if user.Result[characterName] > user.Result[resultName] {
			resultName = characterName
		}
	}
	reply := api.NewMessage(id, "Вы " + resultName)
	reply.ReplyMarkup = api.NewInlineKeyboardMarkup(
		api.NewInlineKeyboardRow(
			api.NewInlineKeyboardButtonData("Пройти тест заново", "start"),
		),
	)
	botUser.Send(reply)
}

func onInlineData(callback *api.CallbackQuery) {
	chatId := callback.Message.Chat.ID
	botUser.Request(api.NewDeleteMessage(chatId, callback.Message.MessageID))
	user, contains := users[chatId]
	if callback.Data == "start" || !contains {
		StartTest(chatId)
		return
	}
	//user := users[chatId]
	user.Result[callback.Data] += 1
	user.CurrentAnswer += 1
	if user.CurrentAnswer == len(questions) {
		PrintResult(chatId)
	} else {
		SendQuestion(chatId, user.CurrentAnswer)
	}
}

func onNewMessage(msg *api.Message) {
	switch msg.Text {
	case "/start":
		StartTest(msg.Chat.ID)
		break
	}
}

func StartTest(id int64) {
	users[id] = &user{
		CurrentAnswer: 0,
	}
	user := users[id]
	user.Result = make(map[string]int)
	for characterId := range characters {
		user.Result[characters[characterId]] = 0
	}
	SendQuestion(id, 0)
}

func SendQuestion(id int64, num int) {
	reply := api.NewMessage(id, "")
	var buttons [][]api.InlineKeyboardButton
	questionName := questions[num]
	characterId := 0
	for variable := range questionsData[questionName] {
		button := api.InlineKeyboardButton{
			Text: questionsData[questionName][variable],
			CallbackData: &characters[characterId],
		}
		buttons = append(buttons, api.NewInlineKeyboardRow(button))
		characterId++
	}
	reply.ReplyMarkup = api.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
	reply.Text = questionName
	botUser.Send(reply)
}