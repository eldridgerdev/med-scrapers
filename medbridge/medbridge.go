package medbridge

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Content []Content `json:"content"`
}

type Content struct {
	Title             string        `json:"title"`
	Brief_description string        `json:"brief_description"`
	Hero_image        string        `json:"hero_image"`
	Instructors       []Instructors `json:"instructors"`
}

type Instructors struct {
	Name_without_ending string `json:"name_without_ending"`
}

type Course struct {
	Url, Image, Title, Presenters string
}

func ScrapeMedbridge() {
	resp, err := http.Get("https://www.medbridge.com/api/v3/courses/filter?&accreditation_state=1&accreditation_discipline=1&sort_by=approved&discipline_id=1")
	if err != nil {
		log.Fatalln("Error retrieving MedBridge courses")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error parsing MedBridge courses")
	}
	var response Response

	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatalf("Unable to marshal jsone due to %s", jsonErr)
	}

	marshalData, jsonMarshalErr := json.MarshalIndent(response, "", "  ")
	if jsonMarshalErr != nil {
		log.Fatalf("Unable to marshal jsone due to %s", jsonMarshalErr)
	}
	os.WriteFile("medbridge-courses.json", marshalData, 0644)
	fmt.Println("MedBridge Courses parsed!")
}

func main() {
	ScrapeMedbridge()
}
