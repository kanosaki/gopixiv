package pixiv

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/k0kubun/pp"
)

var (
	DUMP_COMMUNICATION_DIR = flag.String("dump", "", "dump all requests as file")
	DEBUG                  = flag.Bool("debug", false, "debug mode (verbose output)")
)

func main() {
	flag.Parse()
	if *DEBUG {
		logrus.SetLevel(logrus.DebugLevel)
	}
	pixiv, err := NewFromConfigFile("private/config.json")
	if err != nil {
		panic(err)
	}
	ranking, err := pixiv.Ranking("all", "daily", 50, nil, 1)
	if err != nil {
		panic(err)
	}
	for _, item := range ranking {
		if item.Work.PageCount > 1 {
			pp.Println(item)
			fullFirst, err := item.Work.Fill(pixiv)
			if err != nil {
				panic(err)
			}
			pp.Println(fullFirst)
			break
		}
	}
	//first := ranking[0].Work
	//search, err := pixiv.Search("チノ", 1)
	//if err != nil {
	//	panic(err)
	//}
	//first := search[0]
	//pp.Println(first)
	//img, err := first.OpenImage(pixiv, SIZE_LARGE, 1)
	//if err != nil {
	//	panic(err)
	//}
	//defer img.Close()
	//dst, err := os.Create(fmt.Sprintf("/Users/saki/table/sample.%s", first.Extension()))
	//if err != nil {
	//	panic(err)
	//}
	//defer dst.Close()
	//io.Copy(dst, img)
}
