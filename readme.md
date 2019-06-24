# Check if records with given `ObjectID`s exist in Algolia

> [Algolia][1] is a search engine service provider

## Usage

```sh
$ go run main.go -app_id $_APP_ID -api_key $_READ_KEY -index "index_name" -input all_ids

$ go run main-chunked.go -app_id $_APP_ID -api_key $_READ_KEY -index "index_name" -input all_ids

$ go run check-attribute.go -app_id $_APP_ID -api_key $_READ_KEY -index "index_name" -input all_ids -attr name -show_attr
```

[1]: https://www.algolia.com
