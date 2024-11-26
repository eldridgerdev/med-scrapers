package askai

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
		Data []struct {
			Id         int      `json:"id"`
			Categories []string `json:"categories"`
		} `json:"data"`
	}
	GroqResponse struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
		Choices []struct {
			Id      string `json:"id"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			Index int `json:"index"`
		} `json:"choices"`
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
	Model          string               `json:"model"`
	Stop           any                  `json:"stop"`
	ResponseFormat GroqResponseFormat   `json:"response_format"`
	Messages       []GroqRequestMessage `json:"messages"`
	Stream         bool                 `json:"stream"`
	Temperature    int                  `json:"temperature"`
	MaxTokens      int                  `json:"max_tokens"`
	TopP           int                  `json:"top_p"`
}

type AICallProps struct {
	Data   *any
	Prompt string
}

func AiCall(props AICallProps) FinalResult {
	viper.SetConfigFile("./.env")
	viper.ReadInConfig()
	client := &http.Client{}
	groqMessages := []GroqRequestMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant",
		}, {
			Role:    "user",
			Content: props.Prompt,
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

	var b GroqResponse
	jsonerr := json.Unmarshal(resBody, &b)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	if b.Error != nil {
		log.Fatal(b.Error.Message)
	}

	var finalRes FinalResult
	finalerr := json.Unmarshal([]byte(b.Choices[0].Message.Content), &finalRes)
	if finalerr != nil {
		fmt.Println("ERROR ERROR ERROR")
		fmt.Println(finalRes)
		log.Fatal(finalerr)
	}

	fmt.Println(finalRes)
	return finalRes
}

func main() {
	p := AICallProps{
		Prompt: "Give me an example of user data in JSON format without newlines",
	}
	AiCall(p)
}
