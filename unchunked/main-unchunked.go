package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

func main() {
	fmt.Println("go run main.go -app_id $APP_ID -api_key $READ_KEY -index \"index_name\" -input test_file")
	input := flag.String("input", "", "the file to read the IDs from")
	app_id := flag.String("app_id", "", "application ID for the algolia cluster to use")
	api_key := flag.String("api_key", "", "API key for the algolia cluster being used")
	indexName := flag.String("index", "", "Name of the index this check is being performed on")

	flag.Parse()

	client := algoliasearch.NewClient(*app_id, *api_key)
	index := client.InitIndex(*indexName)

	b, err := ioutil.ReadFile(*input)
	if err != nil {
		log.Fatal(err)
	}

	listOfObjectsOrig := strings.Split(string(b), "\n")

	listOfObjects := []string{}

	for _, v := range listOfObjectsOrig {
		if len(v) == 0 {
			continue
		}
		listOfObjects = append(listOfObjects, v)
	}

	const algoliaMaxGetSize = 1000
	objects, err := index.GetObjects(listOfObjects)
	if err != nil {
		log.Fatal(err)
	}

	exists := 0
	for _, v := range objects {
		if v != nil {
			exists += 1
		}
	}

	log.Printf("Return %d; %d/%d exist!",
		len(objects),
		exists,
		len(listOfObjects),
	)
}
