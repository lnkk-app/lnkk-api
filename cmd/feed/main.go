package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()

	feed, _ := fp.ParseURL("https://opensource.com/feed")
	//feed, _ := fp.ParseURL("https://feeds.metaebene.me/cre/m4a")
	// feed, _ := fp.ParseURL("http://oncloud.deloitte.libsynpro.com/rss")

	fmt.Println(feed)
}
