package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/scrape"
)

func main() {
	parse, err := scrape.Parse(
		context.Background(),
		"https://habr.com/ru/companies/bothub/articles/804077/",
	)
	if err != nil {
		log.Fatal(err)
	}

	encoded, err := json.MarshalIndent(parse, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", string(encoded))
}
