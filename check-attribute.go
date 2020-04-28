package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

func main() {
	fmt.Println("go run main-chunked.go -app_id $APP_ID -api_key $READ_KEY -index \"index_name\" -input test_file -attr attr_name -showAttr")
	input := flag.String("input", "", "the file to read the IDs from")
	app_id := flag.String("app_id", "", "application ID for the algolia cluster to use")
	api_key := flag.String("api_key", "", "API key for the algolia cluster being used")
	indexName := flag.String("index", "", "Name of the index this check is being performed on")
	attributeName := flag.String("attr", "", "Name of the attribute whose value you want to check")
	showAttr := flag.Bool("show_attr", false, "Whether to show the attribute that we fetched or not")

	flag.Parse()

	if *attributeName == "" {
		flag.PrintDefaults()
		log.Fatal(fmt.Errorf("Attribute name MUST be non-empty. Use check-record-single-index.go if you are interested only in existence of objects"))
	}

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
	numChunks := int(math.Ceil(float64(len(listOfObjects)) / float64(algoliaMaxGetSize)))
	exists := 0
	attrExists := 0

	for i := 0; i < numChunks; i++ {
		startIndex := i * algoliaMaxGetSize

		endIndex := (i + 1) * algoliaMaxGetSize
		if len(listOfObjects) < endIndex {
			endIndex = len(listOfObjects)
		}

		thisChunk := listOfObjects[startIndex:endIndex]

		thisChunkExists := 0

		objects, err := index.GetObjects(thisChunk)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range objects {
			if v != nil {
				thisChunkExists += 1
			}
		}

		exists += thisChunkExists

		thisChunkAttrExists := 0
		for _, v := range objects {
			if v != nil {
				attr := v[*attributeName]
				if attr != nil {
					thisChunkAttrExists += 1
				}
				if *showAttr {
					fmtString := "%s: %v"

					switch attr.(type) {
					case float64:
						fmtString = "%s: %.2f"
					}

					log.Printf(fmtString, v["objectID"], attr)
				}
			}
		}

		attrExists += thisChunkAttrExists

		log.Printf("Chunk %d finished: Returned: %d; %d/%d exist! Attribute: %d/%d;",
			i+1,
			len(objects),
			thisChunkExists,
			len(thisChunk),
			thisChunkAttrExists,
			len(thisChunk),
		)
	}

	log.Printf("Objects: %d/%d exist! Attribute %s: %d/%d",
		exists,
		len(listOfObjects),
		*attributeName,
		attrExists,
		len(listOfObjects),
	)
}
