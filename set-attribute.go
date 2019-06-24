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
	fmt.Println("go run check-attribute.go -app_id $ALGOLIA_PROD_APP_ID -api_key $ALGOLIA_PROD_READ_KEY -index \"items.20170801.sort_by_lowest_price\" -input test_file -attr promote_expired_time -set_nil")
	input := flag.String("input", "", "the file to read the IDs from")
	app_id := flag.String("app_id", "", "application ID for the algolia cluster to use")
	api_key := flag.String("api_key", "", "API key for the algolia cluster being used")
	indexName := flag.String("index", "", "Name of the index this check is being performed on")
	attributeName := flag.String("attr", "", "Name of the attribute whose value you want to check")
	setNil := flag.Bool("set_nil", false, "Whether to set the value of this attribute in these items to nil")

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
		if v == "" {
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

	objectsWithAttr := []algoliasearch.Object{}

	attrExists := 0
	for _, v := range objects {
		if v != nil {
			if v[*attributeName] != nil {
				attrExists += 1
				objectsWithAttr = append(objectsWithAttr, v)
			}
		}
	}

	log.Printf("%d/%d have %s!",
		attrExists,
		len(listOfObjects),
		*attributeName,
	)

	if *setNil {
		objects := []algoliasearch.Object{}
		for _, obj := range objectsWithAttr {
			id, _ := obj.ObjectID()
			if id != "" {
				objects = append(objects, algoliasearch.Object{
					"objectID":     id,
					*attributeName: int32(1),
				})
			}
		}
		res, err := index.PartialUpdateObjects(objects)
		if err != nil {
			log.Printf("ERROR: Error while partially updating objects: %v", err)
			return
		}

		log.Print("Partially updated objects: %v", res)
		log.Print("Payload: %v", objects)
	}
}
