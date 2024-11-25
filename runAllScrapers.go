package main

import (
	"encoding/json"
	"fmt"
	"log"
	askai "medCourseFinder/askAI"
	"medCourseFinder/medbridge"
	"os"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./.env")
	viper.ReadInConfig()
	err := os.Remove("./medbridge-courses.json")
	medbridge.ScrapeMedbridge(medbridge.MedBridgeOpts{Limit: 1})
	medbridgeData, err := os.ReadFile("./medbridge-courses.json")
	if err != nil {
		log.Fatal(err)
	}
	var jsondata medbridge.Response
	jsonerr := json.Unmarshal(medbridgeData, &jsondata)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	jsonstring, err := json.Marshal(jsondata)
	prompt := "Using this JSON data: " + string(jsonstring) + "  Add a 'Category' array to each content item"
	fmt.Println(prompt)
	p := askai.AICallProps{
		Prompt: prompt,
	}
	data := askai.AiCall(p)
	fmt.Println(data)
}
