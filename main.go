package main

import (
	"github.com/k0kubun/pp"
	"fmt"
	"io"
	"os"
	"flag"
	"github.com/Sirupsen/logrus"
)

var (
	DUMP_COMMUNICATION_DIR = flag.String("dump", "", "dump all requests as file")
	DEBUG = flag.Bool("debug", false, "debug mode (verbose output)")
)

func init() {
	flag.Parse()
	if *DEBUG {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	pixiv, err := NewFromConfigFile("private/config.json")
	if err != nil {
		panic(err)
	}
	ranking, err := pixiv.Ranking("all", "daily", 50, nil, 1)
	if err != nil {
		panic(err)
	}
	for _, item := range ranking {
		if item.Work.ItemType == ITEM_TYPE_MANGA {
			pp.Println(item)
		}
	}
	first := ranking[0].Work
	//search, err := pixiv.Search("チノ", 1)
	//if err != nil {
	//	panic(err)
	//}
	//first := search[0]
	pp.Println(first)
	img, err := first.OpenImage(pixiv, SIZE_LARGE, 1)
	if err != nil {
		panic(err)
	}
	defer img.Close()
	dst, err := os.Create(fmt.Sprintf("/Users/saki/table/sample.%s", first.ContentType))
	if err != nil {
		panic(err)
	}
	defer dst.Close()
	io.Copy(dst, img)
}
