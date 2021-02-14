package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
)

var cache struct {
	URL     string `json:"url"`
	Content int    `json:"content"`
}

var (
	feedURL   = flag.String("f", "", "The RSS feed you wish to parse from")
	feedView  = flag.Int("v", -1, "The RSS feed entry you wish to view")
	cachePath = ".cache.json"
)

func main() {

	flag.Parse()
	if *feedView <= 0 && len(*feedURL) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: goread [-f feed url] ...\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *feedURL != "" {
		readFeed()
	}

	if *feedView > 0 {
		showContent()
	}

}

func readFeed() {
	file, err := os.Create(cachePath)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
		return
	}
	defer file.Close()

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(*feedURL)

	if feed.Author != nil {
		fmt.Println(feed.Title)
		fmt.Println(feed.Author.Name)
	} else {
		fmt.Println(feed.Title)
	}

	contentLen := len(feed.Items)
	for i := 0; i < contentLen; i++ {
		contentTitle := feed.Items[i].Title
		formattedTitles := fmt.Sprintf("[%d] %s", i+1, contentTitle)
		fmt.Println(formattedTitles)
	}

	// writing cache file to store how many options to read there were
	// Also storing the url for future use
	cachejson := fmt.Sprintf("{\"url\":\"%s\", \"content\": %d}", *feedURL, contentLen)

	var cacheBlob = []byte(cachejson)

	err = json.Unmarshal(cacheBlob, &cache)
	if err != nil {
		log.Fatal(err.Error())
	}

	cacheJSONBlob, _ := json.Marshal(cache)
	err = ioutil.WriteFile(cachePath, cacheJSONBlob, 0644)
}

func showContent() {
	dealersChoice := *feedView
	configFile, err := os.Open(cachePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&cache); err != nil {
		log.Fatalln(err.Error())
	}

	if dealersChoice > cache.Content {
		log.Fatalf("Error choice is outside of bounds number must be no greater than %d", cache.Content)
	}

	fmt.Printf("%s \n%d\n", cache.URL, cache.Content)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(cache.URL)

	text := fmt.Sprintf("%s", feed.Items[dealersChoice-1])
	strippedContent := strip.StripTags(text)

	if feed.Author != nil {
		fmt.Println(feed.Title)
		fmt.Println(feed.Author.Name)
	} else {
		fmt.Println(feed.Title)
	}
	fmt.Println(strippedContent)
}
