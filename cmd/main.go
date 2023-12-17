package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	hg "com.sander/hugging-face-api"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Bot token and username
	botToken := "BOT_TOKEN"

	// Create a Telegram bot API instance
	bot, err := telegram.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Handle incoming updates
	updates, _ := bot.GetUpdatesChan(telegram.UpdateConfig{})

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Photo != nil {
			file, err := bot.GetFile(telegram.FileConfig{FileID: ((*update.Message.Photo)[0]).FileID})
			if err != nil {
				log.Println(err)
				continue
			}

			fileData, err := http.Get(file.Link(botToken))
			if err != nil {
				log.Println(err)
				continue
			}

			b, _ := io.ReadAll(fileData.Body)

			rs, err := GetText(b)
			if err != nil {
				_, err = bot.Send(telegram.MessageConfig{Text: err.Error(), BaseChat: telegram.BaseChat{ChatID: update.Message.Chat.ID}})
				continue
			}

			var response string
			for _, r := range rs {
				if len(r.GeneratedText) > 0 {
					response += r.GeneratedText
				} else {
					response += fmt.Sprintf("%v: %v\n", r.Score, r.Label)
				}
			}

			_, err = bot.Send(telegram.MessageConfig{Text: response, BaseChat: telegram.BaseChat{ChatID: update.Message.Chat.ID}})
			if err != nil {
				log.Println(err)
			}
		}
	}
}

type Response struct {
	Score         float64 `json:"score"`
	Label         string  `json:"label"`
	GeneratedText string  `json:"generated_text"`
}

var TextToImage hg.ModelContext = hg.ModelContext{
	ModelId: "microsoft/resnet-50",
}

func GetText(file []byte) ([]Response, error) {
	res, _ := TextToImage.Request(file)
	bodyBytes, _ := hg.ReadAllClose(&res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, HttpError{Code: res.StatusCode, Message: string(bodyBytes)}
	}
	var response []Response
	json.Unmarshal(bodyBytes, &response)
	return response, nil
}

type HttpError struct {
	Code    int
	Message string
}

func (he HttpError) Error() string {
	return fmt.Sprintf("Error Code %v\n%v\n", he.Code, he.Message)
}
