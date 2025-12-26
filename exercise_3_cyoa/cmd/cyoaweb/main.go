package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"github.com/bitswright/go_mini_projects/exercise_3_cyoa"
)

func main() {
	storyFilename := flag.String("story", "gopher.json", "JSON file with CYOA story")
	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *storyFilename)

	file, err := os.Open(*storyFilename)
	if err != nil {
		log.Fatalf("Error while reading file %s. Error: %s", *storyFilename, err)
		// Todo: How is panic different from log.Fatalf and when to use which???
		panic(err)
	}

	story, err := cyoa.GetStoryFromJson(file)
	if err != nil {
		log.Fatalf("Error while decoding json. Error: %s", err)
	}



	fmt.Printf("%+v\n", story)
}