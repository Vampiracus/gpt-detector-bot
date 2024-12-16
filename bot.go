package main

import (
	"context"
	"os"
	"os/signal"

	"fmt"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Post struct {
	Prediction float64 `json:"prediction"`
}

var ip string = os.Getenv("SERVICE_IP")
var token string = os.Getenv("BOT_TOKEN")

func main() {
	if token == "" {
		panic("BOT_TOKEN is not set")
	}
	if ip == "" {
		panic("SERVICE_IP is not set")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot started")
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text == "/help" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "Я отправляю за вас запрос в API. Просто введите любой текст\n```\ncurl -X POST \"http://" + ip + ":5000/predict\" \\\n-H \"Content-Type: application/json\" \\\n-d '{\"text\": \"Пример текста для предсказания\"}'\n```",
			ParseMode: "Markdown",
		})
		return
	}

	posturl := "http://" + ip + ":5000/predict"
	body := []byte(`{ "text":  "` + update.Message.Text + `" }`)
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	post := &Post{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		panic(derr)
	}
	fmt.Println("RESPONSE", post.Prediction)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Вероятность: %.5f", post.Prediction),
	})
}
