package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

const (
	GROQ_MODEL = "llama3-8b-8192"
)

type (
	FinalResult struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	GroqResponse struct {
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Id string `json:"id"`
	}
	GroqResponseFormat struct {
		Type string `json:"type"`
	}
	GroqRequestMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

type GroqBody struct {
	ResponseFormat GroqResponseFormat   `json:"response_format"`
	Messages       []GroqRequestMessage `json:"messages"`
	Model          string               `json:"model"`
	Temperature    int                  `json:"temperature"`
	MaxTokens      int                  `json:"max_tokens"`
	TopP           int                  `json:"top_p"`
	Stream         bool                 `json:"stream"`
	Stop           any                  `json:"stop"`
}

func main() {
	aiCall()
}

func aiCall() FinalResult {
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()
	client := &http.Client{}
	groqMessages := []GroqRequestMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant",
		}, {
			Role:    "user",
			Content: "Give me an example of user data in JSON format without newlines",
		},
	}

	data := GroqBody{
		Messages: groqMessages,
		ResponseFormat: GroqResponseFormat{
			Type: "json_object",
		},
		Model:       GROQ_MODEL,
		Temperature: 1,
		MaxTokens:   1024,
		Stop:        nil,
		TopP:        1,
		Stream:      false,
	}

	body, err := json.Marshal(data)
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+viper.GetString("GROQ_API"))

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resBody))
	var b GroqResponse
	jsonerr := json.Unmarshal(resBody, &b)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}

	var finalRes FinalResult
	finalerr := json.Unmarshal([]byte(b.Choices[0].Message.Content), &finalRes)
	if finalerr != nil {
		log.Fatal(err)
	}
	// pretty, err := json.MarshalIndent(b.Choices[0].Message.Content, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(finalRes)
	return finalRes
}
