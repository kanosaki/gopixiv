# gopixiv -- Pixiv API wrapper for Golang

NOTE: This is "unofficial" repository for "private" REST api.  This repository might be unexpectedly closed in future.

NOTE: This repository is under developing.

# License

MIT License

# Usage

```go
pixiv, _ := gopixiv.New("<<<client_id>>>", "<<<client_secret>>>", "<<<username>>>", "<<<password>>>")

// search api
// Simple API:  Search("word", page_number)
itemsFromSimple, err := pixiv.Search("word", 1)

// Query API
query := pixiv.SearchQuery("word")
// you can configure more about search parameter
query.Mode = gopixiv.SEARCH_MODE_CAPTION // default: SEARCH_MODE_TAG
query.Sort = gopixiv.SEARCH_SORT_POPULARITY // note: this option required paid(premium) account. default: SEARCH_SORT_DATE
// etc...
// and exec query
itemsFromQuery, err := query.Get(page_number)



```


# Features

* Image fetching
* Search API
* Ranking API


# Future works

* Versatile data type support (manga, ugoira, etc...)
* More API support
  * User API
  
