package main

import "github.com/k0kubun/pp"

func main() {
	pixiv, err := NewFromConfigFile("private/config.json")
	if err != nil {
		panic(err)
	}
	//ranking, err := pixiv.Ranking("all", "daily", 50, nil, 1)
	//if err != nil {
	//	t.Error(err)
	//}
	//if len(ranking) == 0 {
	//	t.Error("Empty result!")
	//}
	search, err := pixiv.Search("チノ", 1)
	if err != nil {
		panic(err)
	}
	pp.Println(search)
}
