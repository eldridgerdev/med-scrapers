package main

import (
	"encoding/json"
	"fmt"
	"log"
	askai "medCourseFinder/askAI"
	"medCourseFinder/medbridge"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./.env")
	viper.ReadInConfig()
	os.Remove("./medbridge-courses.json")
	os.Remove("./medbridge-course-categories.json")
	medbridge.ScrapeMedbridge(medbridge.MedBridgeOpts{Limit: 15})
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
	var promptBuilder strings.Builder

	promptBuilder.WriteString("Using this JSON data: ")
	promptBuilder.WriteString(string(jsonstring))
	promptBuilder.WriteString("For each content item generate a list of categories in an array that gives each object multiple categories based on body part and physical therapy technique")
	promptBuilder.WriteString("Return the results as a JSON object with an \"id\" that correlates to the content item and a \"categories\" which is an array of the generated category strings")
	promptBuilder.WriteString("Create the json object under a top-level array \"data\"")
	//	promptBuilder.WriteString("Return a JSON object following this example: {id: 4641, categories: ['Physical Therapy', 'Shoulder']} that includes the id of the content item and a new 'category' array that gives each object multiple categories based on body part and physical therapy technique")
	p := askai.AICallProps{
		Prompt: promptBuilder.String(),
	}
	data := askai.AiCall(p)
	marshalData, jsonMarshalErr := json.MarshalIndent(data, "", "  ")
	if jsonMarshalErr != nil {
		log.Fatalf("Unable to marshal jsone due to %s", jsonMarshalErr)
	}

	os.WriteFile("medbridge-course-categories.json", marshalData, 0644)
	fmt.Println("Final Data: ")
	fmt.Println(data)
}
