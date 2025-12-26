package cyoa

import (
	"io"
	"encoding/json"
)

type Story map[string]Chapter

type Chapter struct {
	Title      	string     		`json:"title"`
	Paragraphs 	[]string 		`json:"story"`
	Options    	[]ChapterOption `json:"options"`
}

type ChapterOption struct {
	Text 	string 	`json:"text"`
	Chapter string	`json:"arc"`
}

func GetStoryFromJson(reader io.Reader) (Story, error) {
	// Note:
	// 		In json Marshall and Unmarshall we are supposed to pass []byte
	// 		But in NewDecoder we are suppose to pass io.Reader
	jsonDecoder := json.NewDecoder(reader)
	var story Story
	if err := jsonDecoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}