package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

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

	text := update.Message.Text
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			text = text[:i] + "\\n" + text[i+1:]
			i++
		}
	}

	body := []byte(`{ "text":  "` + text + `" }`)
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("[Request error]", err)
		return
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("[Client error]", err)
		return
	}
	defer res.Body.Close()
	post := &Post{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		fmt.Println("[decoder error]", derr)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("[decoder error]", err)
		}
		fmt.Println(string(body))
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Вероятность: %.5f", post.Prediction),
	})
}
