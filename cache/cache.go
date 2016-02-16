package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/smalltree0/beego_goblog/models/manage"
	"github.com/smalltree0/com/log"
)

type cache struct {
	BackgroundLeftBar  map[string]string
	BackgroundLeftBars []*admin.Leftbar
	AboutContent       string
}

var Cache = NewCache()

func NewCache() *cache {
	return &cache{BackgroundLeftBar: make(map[string]string)}
}

func init() {
	doReadBackLeftBarConfig()
	doAboutContent()
}

func doReadBackLeftBarConfig() {
	path, _ := os.Getwd()
	f, err := os.Open(path + "/conf/backleft.conf")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &Cache.BackgroundLeftBars)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range Cache.BackgroundLeftBars {
		if v.ID != "" {
			Cache.BackgroundLeftBar[v.ID] = v.ID
		}
	}
}

func doAboutContent() {
	path, _ := os.Getwd()
	f, err := os.Open(path + "/conf/about.md")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	Cache.AboutContent = string(b)
}
